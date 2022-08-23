package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"logger-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.LogEntry
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	if err := l.models.LogEntry.Insert(logEntry); err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}
	return &logs.LogResponse{Result: "logged"}, nil
}

func (app *Config) gRPCListen() {
	fmt.Println("grpc listrning on port: ", grpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", grpcPort))
	if err != nil {
		log.Fatal(err)
	}

	defer listen.Close()

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, &LogServer{models: app.Models})

	log.Printf("grpc server started on port: %s", grpcPort)

	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
