package queue

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

func InitRabbitMQ() {
	var err error
	Conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("RabbitMQ connection error: %v", err)
	}

	Channel, err = Conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ channel error: %v", err)
	}

	_, err = Channel.QueueDeclare(
		"click_events", // name
		true,           // durable
		false,          // autoDelete
		false,          // exclusive
		false,          // noWait
		nil,            // args
	)
	if err != nil {
		log.Fatalf("Queue declare error: %v", err)
	}
}
