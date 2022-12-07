package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/googleapis/enterprise-certificate-proxy/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .emv file")
	}

	MongoDb := os.Getenv("MONGODB_URL")

	mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Database")

	return client
}

var Client *mongo.Client = DbInstance()

func OpenCollection(collectionName string) *mongo.Collection {
	return Client.Database("cluster0").Collection(collectionName)
}
