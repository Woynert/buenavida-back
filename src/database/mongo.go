package database

import (
	"context"
	"log"
	"os"
	"errors"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

func CheckConnection () error {
	
	var mc *mongo.Client = GetClient()

	// check connection
	if err := mc.Ping(context.TODO(), readpref.Primary()); err != nil {

		// disconnect
		err = CloseConnection()
		if (err != nil){
			return err
		}

		// reconnect
		var err error
		err = StartConnection()

		if (err != nil){
			return err
		}

		// check again
		// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo?utm_source=godoc#Client.Connect
		// The Client.Ping method can be used to verify
		// that the connection was created successfully.

		if err = mc.Ping(context.TODO(), readpref.Primary()); err != nil {
			return err
		}
	}

	return nil
}
