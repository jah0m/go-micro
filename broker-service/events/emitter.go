package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	// create a channel
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	return declareExcahnge(ch)
}

// Push sends an event to rabbitmq
func (e *Emitter) Push(event string, severity string) error {
	// create a channel
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		"logs_topic", // exchange
		severity,     // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
