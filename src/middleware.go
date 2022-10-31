package main

import (
	token "woynert/buenavida-api/token"
	db "woynert/buenavida-api/database"

	"context"
    "net/http"
    "github.com/gin-gonic/gin"

	// mongo
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

		var mc *mongo.Client = db.GetClient()

		// check connection
		if err := mc.Ping(context.TODO(), readpref.Primary()); err != nil {

			// disconnect
			err = db.CloseConnection ()
			if (err != nil){
				c.AbortWithStatusJSON(http.StatusInternalServerError,
				gin.H{"message": "Internal server error"})
				return
			}

			// reconnect
			var err error
			err = db.StartConnection()

			if (err != nil){
				c.AbortWithStatusJSON(http.StatusInternalServerError,
				gin.H{"message": "Internal server error"})
				return
			}

			// check again
			// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo?utm_source=godoc#Client.Connect
			// The Client.Ping method can be used to verify
			// that the connection was created successfully.

			if err = mc.Ping(context.TODO(), readpref.Primary()); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError,
				gin.H{"message": "Internal server error"})
				return
			}
		}

	}
}

