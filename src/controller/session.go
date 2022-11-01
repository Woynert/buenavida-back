package controller

import (
	token "woynert/buenavida-api/token"
	db "woynert/buenavida-api/database"

	"fmt"
	"context"
    "net/http"
    "github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const HOST string = "localhost";

// hash passwords with bcrypt
// https://dev.to/nwby/how-to-hash-a-password-in-go-4jae

func Signin(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, "")
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}


func Login(c *gin.Context) {

	var err error
	var form LoginForm

	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if err := c.BindJSON(&form); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	// find user with email

	var mc *mongo.Client = db.GetClient()
	var user db.User
	coll := mc.Database("buenavida").Collection("users")

	err = coll.FindOne(
		context.TODO(),
		bson.D{{"email", form.Email}},
	).Decode(&user)

	if err != nil {
		// ErrNoDocuments: the filter did not match any documents
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"message": "Wrong user/password"})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal server error"})
			return
		}
	}

	// check password

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil{
		c.AbortWithStatusJSON(http.StatusUnauthorized,
		gin.H{"message": "Wrong user/password"})
		return
	}

	// generate tokens

	var tokeninfo token.TokenInfo = token.TokenInfo{
		UserID: "0001",
	}

	accessToken, refreshToken, err := token.Create(tokeninfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Could not generate token"})
		return
	}

	c.SetCookie("refreshToken", refreshToken, 10, "/session/refresh", HOST, false, true)
	c.SetCookie("accessToken", accessToken, 10, "/", HOST, false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
}


// get new access token from refresh token
func Refresh (c *gin.Context) {

	// get claims

	claimsAny, exists := c.Get("claims")

	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Invalid token"})
		return
	}

	var claims *token.TokenClaims
	claims = claimsAny.(*token.TokenClaims)
	fmt.Println(claims)
	fmt.Println(claims.EXP)

	// generate tokens

	var tokeninfo token.TokenInfo = token.TokenInfo{
		UserID: claims.UserID,
	}

	accessToken, refreshToken, err := token.Create(tokeninfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Could not generate token"})
		return
	}

	c.SetCookie("refreshToken", refreshToken, 10, "/session/refresh", HOST, false, true)
	c.SetCookie("accessToken", accessToken, 10, "/", HOST, false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully refreshed tokens"})

}

func Logout (c *gin.Context) {

	// delete tokens (cookies)
	c.SetCookie("accessToken", "", -1, "/", HOST, false, true)
	c.SetCookie("refreshToken", "", -1, "/session/refresh", HOST, false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Succesfully logged out"})
}

