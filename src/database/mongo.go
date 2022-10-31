package database

import (
	"context"
	"log"
	"os"
	"errors"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient mongo.Client;

func GetClient () *mongo.Client {
	return &mongoClient
}

func StartConnection () error {

	var err error

	// get MONGODB_URI from enviroment

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		// You must set your 'MONGODB_URI' environmental variable. See
		// https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable 
		return errors.New("MONGODB_URI not set")
	}

	// connect
	newClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	mongoClient = *newClient

	if err != nil {
		return err
	}
	return nil
}

func CloseConnection () error {
	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		return errors.New("Could not disconnect")
	}
	return nil
}

