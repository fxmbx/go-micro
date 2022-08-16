package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",      //type
		true,         //durable?
		false,        //autodeleted?
		false,        //internal?
		false,        //nowaite?
		nil,          //arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //name
		false, //durable? not durable just get rid of it when done with it
		false, //delete when unused? do not auto delete if unused
		true,  //exclusive channel?
		false, //nowait?
		nil,   //arguments?
	)
}
