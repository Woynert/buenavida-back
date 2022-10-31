package main

import (
	controller "woynert/buenavida-api/controller"

    "github.com/gin-gonic/gin"
)

func main() {
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

	router.Run(":8070")

}

