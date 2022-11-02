package main

import (
	controller "woynert/buenavida-api/controller"
	db "woynert/buenavida-api/database"
	"log"

    "github.com/gin-gonic/gin"
)

func main() {

	// mongodb search index
	if err := db.CheckConnection(); err != nil {
		log.Fatal(err)
	}
	db.PopulateNgrams()
	db.CreateIndexNgram()

	// router
	router := gin.Default();

	// session
	router.POST   ("/session/signin" , CheckMongoConnection(), controller.Signin)
	router.POST  ("/session/login"  , CheckMongoConnection(), controller.Login)
	router.DELETE("/session/logout" , CheckAccessToken(), controller.Logout)
	router.GET   ("/session/refresh", CheckRefreshToken(), controller.Refresh)

	// cart
	//router.GET("/cart", CheckMongoConnection(&mongoClient), CheckAccessToken(), controller.CartGetItems)

	// favorite

	// store
	router.GET ("/store", CheckMongoConnection(), controller.StoreFilterItems)

	router.Run(":8070")

}

