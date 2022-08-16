package event

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setUp() error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return declareExchange(ch)
}

func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return nil
	}
	defer channel.Close()

	log.Println("pushing to channel")

	err = channel.PublishWithContext(context.TODO(), "logs_topic", severity, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(event),
	})
	if err != nil {
		return err
	}
	return nil
}
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}
	err := emitter.setUp()
	if err != nil {
		log.Println(err)
		return Emitter{}, nil
	}
	return emitter, nil
}
