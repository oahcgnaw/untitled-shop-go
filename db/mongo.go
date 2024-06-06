package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://gdd:Towski2997@cluster0.h3aqe8x.mongodb.net")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

    err = client.Ping(ctx, nil)

    if err != nil {
        log.Fatal(err)
    }

	log.Println("Connected to MongoDB!")
	Client = client
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("eshop_dev").Collection(collectionName)
}