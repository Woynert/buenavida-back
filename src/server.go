package main

import (
	controller "woynert/buenavida-api/controller"

    "github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default();

	// session
	router.GET   ("/session/signin" , controller.Signin)
	router.POST  ("/session/login"  , controller.Login)
	router.DELETE("/session/logout" , CheckAccessToken(), controller.Logout)
	router.GET   ("/session/refresh", CheckRefreshToken(), controller.Refresh)

	// cart
	router.GET("/cart", CheckAccessToken(), controller.CartGetItems)

	// favorite

	// store

	router.Run(":8070")

}

