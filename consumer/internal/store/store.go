package store

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoStore(uri, dbName string) *MongoStore {
	var client *mongo.Client
	var err error

	const (
		maxRetries    = 10
		retryInterval = 5
	)

	for i := 0; i < maxRetries; i++ {
		client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), retryInterval*time.Second)
			defer cancel()

			err = client.Ping(ctx, nil)
			if err == nil {
				log.Println("Connected successfully to MongoDB!")
				return &MongoStore{
					Client:   client,
					Database: client.Database(dbName),
				}
			}
		}

		log.Printf("Error connecting to MongoDB: %v. Trying again in %d seconds...", err, retryInterval)
		time.Sleep(retryInterval * time.Second)
	}

	log.Fatalf("Failed to connect to MongoDB: %v", err)
	return &MongoStore{}
}

func (mongoStore *MongoStore) Close() {
	if err := mongoStore.Client.Disconnect(context.Background()); err != nil {
		log.Fatalf("Failed to disconnect MongoDB: %v", err)
	}
	log.Println("MongoDB connection closed")
}
