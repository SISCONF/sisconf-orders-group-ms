package jobs

import (
	"os"

	"github.com/SISCONF/sisconf-orders-group-ms.git/internal/utils"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewSisconfAmqpBrokerChannel() (*amqp.Connection, *amqp.Channel, <-chan amqp.Delivery) {
	err := godotenv.Load()
	utils.PanicOnError("Couldn't load .env", err)

	rabbitMqHost := os.Getenv("RABBIT_MQ_HOSTNAME")
	connection, err := amqp.Dial(rabbitMqHost)
	utils.PanicOnError("Couldn't dial RabbitMQ", err)

	channel, err := connection.Channel()
	utils.PanicOnError("Failed to open RabbitMQ channel", err)

	queueName := os.Getenv("RABBIT_MQ_ORDERS_GROUP_SHEET_QUEUE_NAME")
	queue, err := channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.PanicOnError("Couldn't declare a queue", err)

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.PanicOnError("Failed to set QoS", err)

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	utils.PanicOnError("Couldn't consume queue", err)

	return connection, channel, msgs
}
