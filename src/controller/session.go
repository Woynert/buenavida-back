package controller

import (
	token "woynert/buenavida-api/token"
	db "woynert/buenavida-api/database"

	"context"
    "net/http"
	"net/mail"
	"fmt"
	"unicode"
    "github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const HOST string = "localhost";

// hash passwords with bcrypt
// https://dev.to/nwby/how-to-hash-a-password-in-go-4jae

type SigninForm struct {
	Firstname string               `json:"firstname"`
	Lastname  string               `json:"lastname"`
	Email     string               `json:"email"`
	Password  string               `json:"password"`
}

func Signin(c *gin.Context) {
	var err error
	var form SigninForm

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


	if !(err != nil) {
		if !(err == mongo.ErrNoDocuments) {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"message": "User already exists"})
			return
		} 
	}

	if form.Firstname == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Name cannot be empty"})
		return
	}

	if form.Lastname == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Lastname cannot be empty"})
		return
	}

	_, err = mail.ParseAddress(form.Email)

	if form.Email == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Email cannot be empty"})
		return
	} else if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Email is invalid"})
        return
	}

	if form.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Password cannot be empty"})
		return
	}

	if len(form.Password) >= 8{
		message := validPassword(form.Password)
		if message != nil{
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": message.Error()})
            return
		}
	}else{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Password is too short"})
        return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(form.Password), 8)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{"message": "Internal server error"})
		return
	}

	data := map[string]interface{}{
		"firstname":form.Firstname,
		"lastname":form.Lastname,
		"email":form.Email,
		"password":string(hashed),
		"favorites":[]string{},
	}

	result, err := coll.InsertOne(context.TODO(),data)

	c.IndentedJSON(http.StatusOK, gin.H{"message": result.InsertedID})
	
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
	c.IndentedJSON(http.StatusOK, "Change me")
}

func Logout (c *gin.Context) {

	// delete tokens (cookies)
	c.SetCookie("accessToken", "", -1, "/", HOST, false, true)
	c.SetCookie("refreshToken", "", -1, "/session/refresh", HOST, false, true)
	c.IndentedJSON(http.StatusOK, "Succesfully logout")
}

func validPassword(s string) error {
	next:
			for name, classes := range map[string][]*unicode.RangeTable{
					"upper case": {unicode.Upper, unicode.Title},
					"lower case": {unicode.Lower},
					"numeric":    {unicode.Number, unicode.Digit},
					"special":    {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
			} {
					for _, r := range s {
							if unicode.IsOneOf(classes, r) {
									continue next
							}
					}
					return fmt.Errorf("password must have at least one %s character", name)
			}
			return nil
}