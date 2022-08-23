package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	port     = "80"
	rpcPort  = "5001"
	mongoUrl = "mongodb://mongo:27017"
	grpcPort = "50001"
)

type Config struct {
	Models data.Models
}

var client *mongo.Client

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	//create conext to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	// infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	// errorLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := Config{
		Models: data.New(client),
	}

	//start web server
	// go app.serve()

	//register the rpc server to tell the app that we'll be accepting rpc request
	err = rpc.Register(new(RPCServer))
	//start rpc server
	go app.rpcListen()

	//LISTEN FOR GRPC CONNECTIONS

	go app.gRPCListen()

	log.Println("Starting service on port :", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// username := os.Getenv("Username")
	// password := os.Getenv("Password")
	//create connectioin options
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connection 😞 : ", err)
		return nil, err
	}
	log.Println("Connected to mongo 😉")
	return c, nil
}
