package controller

import (
	db "woynert/buenavida-api/database"

	"os"
	"fmt"
	"context"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserInformation(c *gin.Context) {
	var err error

	// get userid

	userIdAny, exists := c.Get("userid")

	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Invalid token"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userIdAny.(string))

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Internal server error"})
		return
	}

	// get user if exists

	var mc *mongo.Client = db.MongoGetClient()
	var user db.User
	coll := mc.Database("buenavida").Collection("users")

	err = coll.FindOne(
		context.TODO(),
		bson.D{{"_id", objID}},
	).Decode(&user)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Internal server error"})
		return
	}

	// get favorites info
	// TODO: this can be optimized by using aggregation

	hostIS := ""
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	} else {
		hostIS = os.Getenv("HOST_PRODUCT_IMAGES")
		if hostIS == "" {
			fmt.Println("HOST_PRODUCT_IMAGES not set. Continuing anyway")
		}
	}

	var arrayProduct []db.Product

	for _, x := range user.Favorites {

		// get product if exists

		var product db.Product
		coll := mc.Database("buenavida").Collection("products")
		err = coll.FindOne(
			context.TODO(),
			bson.D{{"_id", x}},
		).Decode(&product)

		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal server error"})
			return
		}

		if hostIS != "" {
			product.ImageUrl = fmt.Sprintf("%s/%s", hostIS, product.ImageUrl)
		}

		arrayProduct = append(arrayProduct, product)
	}

	dataUser := map[string]interface{}{
		"firstname":user.Firstname,
		"lastname":user.Lastname,
		"email":user.Email,
		"favorites":arrayProduct,
	}

	c.IndentedJSON(http.StatusOK,dataUser)
}
