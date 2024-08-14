package main

import (
	"context"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Model data.Models
}

// implement the WriteLog method
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	err := l.Model.LogEntry.Insert(data.LogEntry{
		Name: req.LogEntry.Name,
		Data: req.LogEntry.Data,
	})
	if err != nil {
		return &logs.LogResponse{Result: "Failed"}, nil
	}

	return &logs.LogResponse{Result: "Logged!"}, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Panic(err)
	}

	grpcServer := grpc.NewServer()

	// Register the LogService
	logs.RegisterLogServiceServer(grpcServer, &LogServer{Model: app.Models})
	log.Println("gRPC server is running on port: " + grpcPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Panic(err)
	}

}
