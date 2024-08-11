package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(p RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:     p.Name,
		Data:     p.Data,
		CreateAt: time.Now(),
	})
	if err != nil {
		log.Println("Failed to insert log entry", err)
		return err
	}

	*resp = "Log entry created by RPC: " + p.Name
	return nil
}
