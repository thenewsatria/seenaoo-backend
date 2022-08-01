package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx, cancel = context.WithCancel(context.Background())
var client *mongo.Client
var err error
var databaseName string

func ConnectDB() {
	credential := options.Credential{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	clientOpts := options.Client().ApplyURI("mongodb://" + os.Getenv("DB_HOSTNAME") + ":" + os.Getenv("DB_PORT")).SetAuth(credential)
	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		cancel()
		log.Fatal("Database Connection Error $s", err)
	}

	databaseName = os.Getenv("DB_NAME")
}

func DisconnectDB() {
	if err := client.Disconnect(ctx); err != nil {
		log.Fatal("Error while disconnecting database $s", err)
	}
}

func PingDB() {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Error while pinging database $s", err)
	}
}

func GetDBContext() context.Context {
	return ctx
}

func CancelDBContext() {
	cancel()
}

func GetDBClient() *mongo.Client {
	return client
}

func UseDB() *mongo.Database {
	return client.Database(databaseName)
}
