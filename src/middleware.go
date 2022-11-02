package main

import (
	token "woynert/buenavida-api/token"
	db "woynert/buenavida-api/database"

	"fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)

func CheckAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		var accessToken string
		accessToken, err := c.Cookie("accessToken")

		// access token cookie not found
		if err != nil{
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "Access token not provided"})
		}

		if token.Validate(accessToken){
			c.Next()
		} else { 
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "Invalid access token"})
		}

	}
}

func CheckRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		var refreshToken string
		refreshToken, err := c.Cookie("refreshToken")

		// refresh token cookie not found
		if err != nil{
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "Refresh token not provided"})
		}

		if token.Validate(refreshToken){
			c.Next()
		} else { 
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "Invalid refresh token"})
		}

	}
}

// ensure there is a connection
func CheckMongoConnection() gin.HandlerFunc {
	return func(c *gin.Context){

		err := db.CheckConnection()
		if (err != nil){
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal server error"})
			return
		}

	}
}

