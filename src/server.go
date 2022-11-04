package main

import (
	controller "woynert/buenavida-api/controller"
	db "woynert/buenavida-api/database"
	"log"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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
	router.Use(cors.Default())

	// image server
	router.StaticFS("/assets", http.Dir("../data/assets"))

	// session
	router.POST  ("/session/signin" , CheckMongoConnection(), controller.Signin)
	router.POST  ("/session/login"  , CheckMongoConnection(), controller.Login)
	router.DELETE("/session/logout" , CheckAccessToken(), controller.Logout)
	router.GET   ("/session/refresh", CheckRefreshToken(), controller.Refresh)

	//favorite
	router.PUT   ("/favorite", CheckMongoConnection(), CheckAccessToken(), controller.AddFavorites)
	router.DELETE("/favorite", CheckMongoConnection(), CheckAccessToken(), controller.RemoveFavorites)

	// user

	router.GET("/user/info", CheckMongoConnection(), CheckAccessToken(), controller.GetUserInformation)

	// cart
	router.POST("/payment", CheckMongoConnection(),  CheckPostgresConnection(), CheckAccessToken(), controller.Payment)

	// store
	router.GET ("/store", CheckMongoConnection(), controller.StoreFilterItems)

	router.Run(":8070")

}

