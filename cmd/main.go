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
		go worker(i, msgs, &wg)
	}

	wg.Wait()
}

func worker(id int, msgs <-chan amqp.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()
	var logMsg string

	for delivery := range msgs {
		log.Printf("Worker %d Received a task\n", id)

		var ordersGroup sisconf.OrdersGroup
		err := json.Unmarshal(delivery.Body, &ordersGroup)
		if err != nil {
			logMsg = fmt.Sprintf("Worker %d: Couldn't read message body: %s", id, err.Error())
			log.Println(logMsg)
			delivery.Nack(false, true)
			continue
		}

		err = files.CreateOrdersGroupXlsxFile(ordersGroup)
		if err != nil {
			logMsg = fmt.Sprintf("Worker %d: Couldn't create spreadsheet: %s", id, err.Error())
			log.Println(logMsg)
			delivery.Nack(false, true)
			continue
		}

		delivery.Ack(false)
		log.Printf("Worker %d: Task acknowledged\n", id)
	}
}
