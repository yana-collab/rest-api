package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitializeMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("could not connect to MongoDB: %w", err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	MongoClient = client
	return client, nil
}

func EnsureCollectionAndDocumentExists(client *mongo.Client, log *slog.Logger) {
	collection := client.Database("test").Collection("collection")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Define the document to be checked/inserted test
	document := bson.D{
		{"id", 1},
		{"url", "https://example.com"},
		{"alias", "example_alias"},
	}

	var result bson.M
	err := collection.FindOne(ctx, bson.M{"id": 1}).Decode(&result)
	if err != nil && err.Error() == "mongo: no documents in result" {
		// If the document doesn't exist, insert it
		_, insertErr := collection.InsertOne(ctx, document, options.InsertOne())
		if insertErr != nil {
			log.Error("Failed to insert document", slog.String("error", insertErr.Error()))
			return
		}
		log.Info("Document successfully inserted into MongoDB collection")
	} else if err != nil {
		log.Error("Failed to find document", slog.String("error", err.Error()))
		return
	} else {
		log.Info("Document already exists in MongoDB collection")
	}
}
