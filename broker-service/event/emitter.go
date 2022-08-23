package event

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

//this emitter struct houses the amqp connect and emitter events
type Emitter struct {
	connection *amqp.Connection
}

//sets up a channet connection and declares an exchange channel
func (e *Emitter) setUp() error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return declareExchange(ch)
}

//publishes a message to the channel with context, exchange, key, etc.
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

//creates a new emitter by accessing the setup event and return an emitter or possible and error
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
