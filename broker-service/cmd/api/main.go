package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const port = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	//connect to rabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()
	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s 😏\n", port)

	//define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	//start server
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

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
