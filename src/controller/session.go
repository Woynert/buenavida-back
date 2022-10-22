package controller

import (
	token "woynert/buenavida-api/token"

	"fmt"
    "net/http"
    "github.com/gin-gonic/gin"
)

const HOST string = "localhost";

func Signin(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, "")
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var form LoginForm

	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if err := c.BindJSON(&form); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	fmt.Println(form.Email)
	fmt.Println(form.Password)

	if form.Email == "mario@mushrom.kingdom" && form.Password == "superpass" {

		var tokeninfo token.TokenInfo = token.TokenInfo{
			UserID: "0001",
		}
		accessToken, refreshToken, err := token.Create(tokeninfo)

		fmt.Println("tokens ", err)
		fmt.Println(accessToken)
		fmt.Println(refreshToken)

		c.SetCookie("refreshToken", refreshToken, 10, "/session/refresh", HOST, false, true)
		c.SetCookie("accessToken", accessToken, 10, "/", HOST, false, true)

		c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
	} else {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Wrong user/password"})
	}
}


// get new access token from refresh token
func Refresh (c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Change me")
}

func Logout (c *gin.Context) {

	// delete tokens (cookies)
	c.SetCookie("accessToken", "", -1, "/", HOST, false, true)
	c.SetCookie("refreshToken", "", -1, "/session/refresh", HOST, false, true)
	c.IndentedJSON(http.StatusOK, "Succesfully logout")
}

