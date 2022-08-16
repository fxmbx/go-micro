package main

import (
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	port = "80"
)

type Config struct {
}

func main() {

	//connect to rabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()

	//start listening for message
	log.Println("Listening for and consuming rabbitmq messages")

	//create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Panic(err)
	}

	//watch queuq and consume event

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

// app := &Config{}

// srv := &http.Server{
// 	Addr:    fmt.Sprintf(":%s", port),
// 	Handler: app.routes(),
// }
// err := srv.ListenAndServe()
// if err != nil {
// 	log.Panic(err)
// }
// }

func connect() (*amqp.Connection, error) {
	var counts int64
	backOff := 1 * time.Second
	var connection *amqp.Connection

	//dont connect until rabbit is ready

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Println("RabbitMQ not yet ready 🐇")
			counts++
		} else {
			connection = c
			log.Println("Connected to rabbitmq 🐰")
			break
		}

		if counts > 5 {
			log.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
