package config

import (
    "context"
    "time"
    "log"
    "fmt"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
    Client *mongo.Client
    DB *mongo.Database
)

func MakeConnection(){
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    uri := "mongodb://Mccall:password@localhost:27017/?authSource=admin"
    clientOptions := options.Client().ApplyURI(uri)

    client , err := mongo.Connect(clientOptions)

    if err != nil {
        log.Fatal("Cannot create MongoDb client: ", err)
    }

    if err = client.Ping(ctx, nil); err != nil {
        log.Fatal("MongoDb ping failed", err)
    }

    fmt.Println("MongoDb connected")

    Client = client
    DB = client.Database("ShortURL_Database")
}
