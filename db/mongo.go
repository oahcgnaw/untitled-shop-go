package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() {
	clientOptions := options.Client().ApplyURI(os.Getenv("DB_URL"))
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
	serve_mode := os.Getenv("SERVE_MODE")
	if serve_mode == "dev" {
		return Client.Database("eshop_dev").Collection(collectionName)
	} else {
		return Client.Database("eshop").Collection(collectionName)
	}
	
}