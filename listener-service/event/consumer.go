package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

///will be used for receiving amqp event from the queue
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setUp()
	if err != nil {
		log.Println("something went wrong: ", err)
		return Consumer{}, err
	}
	return consumer, nil
}

//open up a channel and declare an exchange
func (consumer *Consumer) setUp() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		log.Println(err)
		return err
	}
	return declareExchange(channel)
}

//means of pushing events to rabbitmq
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {

	ch, err := consumer.conn.Channel()
	if err != nil {
		log.Println("Something went wrong ", err)
		return err
	}
	defer ch.Close()

	//declare a random queue and use it
	queue, err := declareRandomQueue(ch)
	if err != nil {
		log.Println("Something went wrong ", err)
		return err
	}
	for _, v := range topics {
		ch.QueueBind(
			queue.Name,
			v,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			log.Println("Something went wrong ", err)
			return err
		}
	}

	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Println("Something went wrong ", err)
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message on [Exchange, Queue] [logs_topic, %s]\n", queue.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		//log what we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		//authenticate
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		log.Println(err)
		return err
	}

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err

	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		// app.errorJson(w, errors.New(""))
		log.Println("something went wrong: 😞 ", err)
		return err
	}

	return nil
}
