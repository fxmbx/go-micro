package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net"
	"net/rpc"
	"time"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

func (rpc *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	// ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo:", err)
		return err
	}

	*resp = "Processed payload via RPC: " + payload.Name
	return nil
}

//starting the rpc server to listen
func (app *Config) rpcListen() error {
	log.Println("starting RPC server on port: ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		rpcCon, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcCon)
	}

}
