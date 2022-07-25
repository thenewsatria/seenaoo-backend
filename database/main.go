package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
var client *mongo.Client
var err error

func ConnectDB(dbUsername, dbPassword, dbHostName, dbPort string) {
	credential := options.Credential{
		Username: dbUsername,
		Password: dbPassword,
	}

	clientOpts := options.Client().ApplyURI("mongodb://" + dbHostName + ":" + dbPort).SetAuth(credential)
	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
}

func DisconnectDB() {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func PingDB() {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
}

func GetDBContext() *context.Context {
	return &ctx
}

func CancelDBContext() {
	cancel()
}

func GetDBClient() *mongo.Client {
	return client
}
