package main

import (
	"fmt"
	event "listener/events"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	conn, err := connect()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ", err)
		os.Exit(1)
	}
	defer conn.Close()

	// start listening for messages
	log.Println("Listening for messages")

	// create a consumer
	consumer, err := event.NewConsumer(conn)
	if err != nil {
		log.Println("Failed to create consumer", err)
		os.Exit(1)
	}
	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println("Failed to listen for messages", err)
		os.Exit(1)
	}
}

// connect to rabbitmq
func connect() (*amqp.Connection, error) {
	var conts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://rabbitmq:password@rabbitmq:5672/")
		if err != nil {
			fmt.Println("RabbitMQ is not ready", err)
			if conts == 5 {
				return nil, err
			}
			conts++
			time.Sleep(backOff)
			backOff = backOff * 2
			continue
		} else {
			connection = c
			log.Println("Connected to RabbitMQ")
			break
		}
	}
	return connection, nil
}
