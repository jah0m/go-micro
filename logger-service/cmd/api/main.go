package main

import (
	"context"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to cancel the connection
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//close the connection
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	//Register the RPC server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Panic(err)
	}
	go app.RPCListen()

	// start the server
	log.Println("Server started on port", webPort)
	// create a new server
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}
	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) RPCListen() {
	log.Println("Starting RPC server on port", rpcPort)
	listen, err := net.Listen("tcp", ":"+rpcPort)
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Panicln("Error connecting", err)
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Panicln("Error connecting", err)
		return nil, err
	}
	log.Println("Connected to MongoDB!")
	return client, nil
}
