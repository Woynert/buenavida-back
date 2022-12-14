package controller

import (
	token "woynert/buenavida-api/token"
	db "woynert/buenavida-api/database"

	"fmt"
	"context"
    "net/http"
	"net/mail"
	"unicode"
    "github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const HOST string = "localhost:8070";
const ACCESS_COOKIE_EXP = 60;
const REFRESH_COOKIE_EXP = 60*60*24;

// hash passwords with bcrypt
// https://dev.to/nwby/how-to-hash-a-password-in-go-4jae

type SigninForm struct {
	Firstname string               `json:"firstname"`
	Lastname  string               `json:"lastname"`
	Email     string               `json:"email"`
	Password  string               `json:"password"`
}

func Ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Authorized"})
}

func Signin(c *gin.Context) {
	var err error
	var form SigninForm

	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if err := c.BindJSON(&form); err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	// find user with email

	var mc *mongo.Client = db.MongoGetClient()
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Name cannot be empty"})
		return
	}

	if form.Lastname == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Lastname cannot be empty"})
		return
	}

	_, err = mail.ParseAddress(form.Email)

	if form.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Email cannot be empty"})
		return
	} else if err != nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Email is invalid"})
        return
	}

	if form.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Password cannot be empty"})
		return
	}

	if len(form.Password) >= 8{
		message := validPassword(form.Password)
		if message != nil{
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message.Error()})
            return
		}
	}else{
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Password is too short"})
        return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(form.Password), 8)

	if err != nil {
		fmt.Println(err)
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

	_, err = coll.InsertOne(context.TODO(),data)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError,gin.H{"message": "Internal server error"})
        return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "User created successfully"})
	
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

	var mc *mongo.Client = db.MongoGetClient()
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
		UserID: user.Id.Hex(),
	}

	fmt.Println(user.Id.Hex())

	accessToken, refreshToken, err := token.Create(tokeninfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Could not generate token"})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refreshToken", refreshToken, REFRESH_COOKIE_EXP, "/session/refresh", HOST, false, true)
	c.SetCookie("accessToken", accessToken, ACCESS_COOKIE_EXP, "/", HOST, false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully logged in"})
}


// get new access token from refresh token
func Refresh (c *gin.Context) {

	// get user id

	userIdAny, exists := c.Get("userid")

	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Invalid token"})
		return
	}

	userId := userIdAny.(string)

	// generate tokens

	var tokeninfo token.TokenInfo = token.TokenInfo{
		UserID: userId,
	}

	accessToken, refreshToken, err := token.Create(tokeninfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Could not generate token"})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refreshToken", refreshToken, REFRESH_COOKIE_EXP, "/session/refresh", HOST, false, true)
	c.SetCookie("accessToken", accessToken, ACCESS_COOKIE_EXP, "/", HOST, false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Successfully refreshed tokens"})

}

func Logout (c *gin.Context) {

	// delete tokens (cookies)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("accessToken", "", -1, "/", HOST, false, true)
	c.SetCookie("refreshToken", "", -1, "/session/refresh", HOST, false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Succesfully logged out"})
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
