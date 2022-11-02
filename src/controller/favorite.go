package controller

import (
	db "woynert/buenavida-api/database"

	"fmt"
	"context"
    "net/http"

    "github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Favorite struct {
	Id         primitive.ObjectID `json:"itemid"`
}

type user struct {
	Favorites  string               `json:"favorites"`
	Email     string               `json:"email"`
}

func AddFavorites(c *gin.Context){
	var err error

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

	var products Favorite

	if err := c.BindJSON(&products); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	var product db.Product
	
	coll = mc.Database("buenavida").Collection("products-search")

	err = coll.FindOne(
		context.TODO(),
		bson.D{{"_id", products.Id}},
	).Decode(&product)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Internal server error"})
		return
	}

	var arrayFavorite []primitive.ObjectID = user.Favorites

	var result bool = false
    for _, x := range arrayFavorite {
        if x == product.Id {
            result = true
            break
        }
    }

	if result {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Product already exists in favorites"})
		return
    }

	arrayFavorite = append(arrayFavorite,product.Id)
	
	
	coll = mc.Database("buenavida").Collection("users")
	_, err = coll.UpdateOne(context.TODO(), 
	bson.D{{"_id", objID}}, 
	bson.D{{"$set", bson.D{{"favorites", arrayFavorite}}}})

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Internal server error"})
	}

	c.IndentedJSON(http.StatusOK,gin.H{"message": "Add Favorite successfully"})
}

func RemoveFavorites(c *gin.Context){

}