package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/files"
	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/jobs"
	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/sisconf"
	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/utils"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

const workerCount int = 5

func main() {
	err := godotenv.Load()
	utils.PanicOnError("Couldn't load .env file", err)

	scheduler := jobs.NewSisConfCronScheduler()
	scheduler.Start()
	log.Println("Successfully started cron scheduler")

	workerQueueName := os.Getenv("RABBIT_MQ_ORDERS_GROUP_SHEET_QUEUE_NAME")
	connection, channel, msgs := jobs.NewSisconfAmqpBrokerChannel()
	defer connection.Close()
	defer channel.Close()
	logMsg := fmt.Sprintf("Waiting for tasks in %s", workerQueueName)
	log.Println(logMsg)

	var wg sync.WaitGroup

	for i := range workerCount {
		wg.Add(1)
		go worker(i, msgs, &wg, channel)
	}

	wg.Wait()
}

func worker(id int, msgs <-chan amqp.Delivery, wg *sync.WaitGroup, channel *amqp.Channel) {
	defer wg.Done()
	const maxRetries = 3
	var logMsg string

	for delivery := range msgs {
		log.Printf("Worker %d Received a task\n", id)

		var ordersGroup sisconf.OrdersGroup
		err := json.Unmarshal(delivery.Body, &ordersGroup)
		if err != nil {
			logMsg = fmt.Sprintf("Worker %d: Couldn't read message body: %s", id, err.Error())
			log.Println(logMsg)
			handleRetry(&delivery, maxRetries, channel)
			continue
		}

		err = files.CreateOrdersGroupXlsxFile(ordersGroup)
		if err != nil {
			logMsg = fmt.Sprintf("Worker %d: Couldn't create spreadsheet: %s", id, err.Error())
			log.Println(logMsg)
			handleRetry(&delivery, maxRetries, channel)
			continue
		}

		delivery.Ack(false)
		log.Printf("Worker %d: Task acknowledged\n", id)
	}
}

func handleRetry(delivery *amqp.Delivery, maxRetries int, channel *amqp.Channel) {
	retries, ok := delivery.Headers["x-retries"].(int32)
	if !ok {
		retries = 0
	}

	if retries < int32(maxRetries) {
		delivery.Headers["x-retries"] = retries + 1

		err := channel.Publish(
			"",
			delivery.RoutingKey,
			false,
			false,
			amqp.Publishing{
				Headers:     delivery.Headers,
				ContentType: "application/json",
				Body:        delivery.Body,
			},
		)

		if err != nil {
			log.Printf("Failed to republish message %s\n", err.Error())
			return
		}

		delivery.Ack(false)
		log.Printf("Message requeued (retry %d of %d)\n", retries+1, maxRetries)
	} else {
		delivery.Nack(false, false)
		log.Printf("Message discarded after %d retries\n", maxRetries)
	}
}
