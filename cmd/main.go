package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/files"
	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/jobs"
	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/sisconf"
	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	utils.PanicOnError("Couldn't load .env file", err)

	scheduler := jobs.NewSisConfCronScheduler()
	scheduler.Start()
	log.Println("Successfully started cron scheduler")

	workerQueueName := os.Getenv("RABBIT_MQ_ORDERS_GROUP_SHEET_QUEUE_NAME")
	msgs := jobs.NewSisconfAmqpBrokerChannel()
	logMsg := fmt.Sprintf("Waiting for tasks in %s", workerQueueName)
	log.Println(logMsg)

	go func() {
		for delivery := range msgs {
			var ordersGroup sisconf.OrdersGroup
			err = json.Unmarshal(delivery.Body, &ordersGroup)
			if err != nil {
				logMsg = fmt.Sprintf("Couldn't read message body: %s", err.Error())
				log.Fatalln(logMsg)
			} else {
				files.CreateOrdersGroupXlsxFile(ordersGroup)
			}
		}
	}()

	var forever chan any

	<-forever
}
