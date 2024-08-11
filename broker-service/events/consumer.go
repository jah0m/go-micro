package event

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
	consumer := &Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		log.Println("Failed to setup consumer", err)
		return nil, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	// create a channel
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	return declareExcahnge(ch)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Consumer) Listen(topic []string) error {
	// create a channel
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	// declare a queue
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topic {
		// bind the queue to the exchange
		err = ch.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			"logs_topic", // exchange
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			return err
		}
	}

	// consume the messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var payload Payload
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Println("Failed to unmarshal message", err)
			}

			// handle payload
			go handlePayload(payload)
		}
	}()

	log.Printf("Waiting for messages [Exchange: logs_topic, Queue: %s, Topics: %v]", q.Name, topic)
	<-forever
	return nil
}

func handlePayload(p Payload) {
	switch p.Name {
	case "logs", "events":
		err := logEvent(p)
		if err != nil {
			log.Println("Failed to log event", err)
		}

	case "auth":
		// authenticate user
	default:
		err := logEvent(p)
		if err != nil {
			log.Println("Failed to log event", err)
		}
	}
}

func logEvent(p Payload) error {
	// create some json send the log micro service
	jsonData, _ := json.MarshalIndent(p, "", "\t")

	// call the service
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//make sure get correct status code
	if response.StatusCode != http.StatusCreated {
		return err
	}

	return nil
}
