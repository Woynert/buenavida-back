package database

import (
	//"context"
	"fmt"
	"log"
	"os"
	"errors"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

var postgresClient *sql.DB

func PgGetClient () *sql.DB {
	return postgresClient
}

func PgStartConnection () error {

    var err error

	// get postgres credentials from enviroment

	if err = godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbuser := os.Getenv("POSTGRES_USER")
	if dbuser == "" {
		return errors.New("POSTGRES_USER not set")
	}
	dbpass := os.Getenv("POSTGRES_PASS")
	if dbpass == "" {
		return errors.New("POSTGRES_PASS not set")
	}

	connStr := "user=" + dbuser +
		" password=" + dbpass +
		" dbname=buenavida" + 
		" sslmode=disable" +
		" host=127.0.0.1" +
		" port=5432" 

    // Get a database handle

    postgresClient, err = sql.Open("postgres", connStr)
    if err != nil {
		fmt.Println(err)
		return err
    }

    err = postgresClient.Ping()
    if err != nil {
		fmt.Println(err)
		return err
    }

	return nil
}

func PgCloseConnection () error {
	if err := postgresClient.Close(); err != nil {
		return errors.New("Could not disconnect")
	}
	return nil
}

func PgCheckConnection () error {
	
	var pg *sql.DB = PgGetClient()

	// check connection

	var err error
	err = nil

	if pg != nil {
		err = pg.Ping()
	}

	if ((pg == nil) || (err != nil)) {

		// disconnect
		if pg != nil {
			err = PgCloseConnection()
			if (err != nil){
				return err
			}
		}

		// reconnect
		var err error
		err = PgStartConnection()

		if (err != nil){
			return err
		}

		// check again
		// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo?utm_source=godoc#Client.Connect
		// The Client.Ping method can be used to verify
		// that the connection was created successfully.

		pg = PgGetClient()
		if err = pg.Ping(); err != nil {
			return err
		}
	}

	return nil
}

