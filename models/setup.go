package models

import(
	"context"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const (
    DbHost     = "localhost"
    DbPort     = "27017"
    DbUser     = "Rachel Protect"
    DbPassword = ""
    DbName     = "my_activity_api"
)

var (
    DB     *mongo.Database
    Client *mongo.Client
)

func ConnectDatabase() {
    uri := "mongodb://"
    if DbUser != "" && DbPassword != "" {
        uri += DbUser + ":" + DbPassword + "@"
    }
    uri += DbHost + ":" + DbPort

    clientOptions := options.Client().ApplyURI(uri)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal("Error connecting to MongoDB:", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Error pinging MongoDB:", err)
    }

    log.Println("Connected to MongoDB!")

    Client = client
    DB = client.Database(DbName)
}

func GetCollection(collectionName string) *mongo.Collection {
    return DB.Collection(collectionName)
}

func DisconnectDatabase() {
    if Client != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        
        if err := Client.Disconnect(ctx); err != nil {
            log.Fatal("Error disconnecting from MongoDB:", err)
        }
        log.Println("Disconnected from MongoDB")
    }
}