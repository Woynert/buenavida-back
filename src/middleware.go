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

		// token cookie not found
		if err != nil{
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "Access token not provided"})
			return
		}

		// validate
		err, claims := token.Validate(accessToken)

		if err != nil{
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "Invalid access token"})
			return
		}

		c.Set("userid", claims.UserID)
	}
}

func CheckRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		var refreshToken string
		refreshToken, err := c.Cookie("refreshToken")

		// refresh token cookie not found
		if err != nil{
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "Refresh token not provided"})
			return
		}

		// validate
		err, claims := token.Validate(refreshToken)

		if err != nil{
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"message": "Invalid refresh token"})
			return
		}

		c.Set("userid", claims.UserID)
	}
}

// ensure there is a mongo connection
func CheckMongoConnection() gin.HandlerFunc {
	return func(c *gin.Context){

		err := db.MongoCheckConnection()
		if (err != nil){
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal server error"})
			return
		}

	}
}

// ensure there is a postgres connection
func CheckPostgresConnection() gin.HandlerFunc {
	return func(c *gin.Context){

		err := db.PgCheckConnection()
		if (err != nil){
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal server error"})
			return
		}

	}
}

