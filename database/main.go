package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx, cancel = context.WithCancel(context.Background())
var client *mongo.Client
var err error
var databaseName string

func ConnectDB(dbUsername, dbPassword, dbHostName, dbPort string) {
	credential := options.Credential{
		Username: dbUsername,
		Password: dbPassword,
	}

	clientOpts := options.Client().ApplyURI("mongodb://" + dbHostName + ":" + dbPort).SetAuth(credential)
	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		cancel()
		log.Fatal("Database Connection Error $s", err)
	}
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

func SetDBName(dbName string) {
	databaseName = dbName
}
