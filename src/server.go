package main

import (
	controller "woynert/buenavida-api/controller"
	db "woynert/buenavida-api/database"
	"log"
	"fmt"

    "github.com/gin-gonic/gin"
	flag "github.com/spf13/pflag"
)

func main() {

	// flag

	var flagCalindex = flag.StringP("calindex", "i", "unset", "Calculate index")
	flag.Lookup("calindex").NoOptDefVal = "set"
	flag.Parse()

	if *flagCalindex != "unset" {

		// mongodb search index
		if err := db.MongoCheckConnection(); err != nil {
			log.Fatal(err)
		}
		db.PopulateNgrams()
		db.CreateIndexNgram()

		fmt.Println("Finish creating index")
		return
	}

	// router
	router := gin.Default();

	// session
	router.GET   ("/session/signin" , controller.Signin)
	router.POST  ("/session/login"  , CheckMongoConnection(), controller.Login)
	router.DELETE("/session/logout" , CheckAccessToken(), controller.Logout)
	router.GET   ("/session/refresh", CheckRefreshToken(), controller.Refresh)

	// cart
	router.POST("/payment", CheckMongoConnection(),  CheckPostgresConnection(), CheckAccessToken(), controller.Payment)

	// favorite

	// store
	router.GET ("/store", CheckMongoConnection(), controller.StoreFilterItems)

	router.Run(":8070")

}

