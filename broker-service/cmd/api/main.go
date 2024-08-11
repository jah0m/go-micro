package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to rabbitmq
	conn, err := connect()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ", err)
		os.Exit(1)
	}
	defer conn.Close()

	app := Config{
		Rabbit: conn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
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
