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

type user struct {
	Favorites  string               `json:"favorites"`
	Email     string               `json:"email"`
}

func AddFavorites(c *gin.Context){
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

	// get product id from Query

	var newFavoriteProductId primitive.ObjectID

	if newFavoriteProductId, err = primitive.ObjectIDFromHex(c.DefaultQuery("itemid", "")); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong ObjectID"})
		return
	}

	// check product exists

	var product db.Product

	coll = mc.Database("buenavida").Collection("products")

	err = coll.FindOne(
		context.TODO(),
		bson.D{{"_id", newFavoriteProductId}},
	).Decode(&product)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
		gin.H{"message": "Product not found"})
		return
	}

	// no repeated product ids

	var arrayFavorite []primitive.ObjectID = user.Favorites

	var result bool = false
	for _, x := range arrayFavorite {
		if x == product.Id {
			result = true
			break
		}
	}

	if result {
		c.AbortWithStatusJSON(http.StatusBadRequest,
		gin.H{"message": "Product already exists in favorites"})
		return
	}

	// finally add

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

	c.IndentedJSON(http.StatusOK,gin.H{"message": "Add Item to Favorite successfully"})
}

func RemoveFavorites(c *gin.Context){
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

	// get product id from Query

	var newFavoriteProductId primitive.ObjectID

	if newFavoriteProductId, err = primitive.ObjectIDFromHex(c.DefaultQuery("itemid", "")); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong ObjectID"})
		return
	}

	// check product exists

	var product db.Product

	coll = mc.Database("buenavida").Collection("products")

	err = coll.FindOne(
		context.TODO(),
		bson.D{{"_id", newFavoriteProductId}},
	).Decode(&product)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Internal server error"})
		return
	}

	// no exists product id in favorites

	var arrayFavorite []primitive.ObjectID = user.Favorites

	var result bool = false
	var index int;
    for i, x := range arrayFavorite {
        if x == product.Id {
            result = true
			index = i;
            break
        }
    }

	if !result {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Product not exists in favorites"})
		return
    }

	// finally remove from favorites

	arrayFavorite = removeIndex(arrayFavorite, index)

	coll = mc.Database("buenavida").Collection("users")
	_, err = coll.UpdateOne(context.TODO(),
	bson.D{{"_id", objID}},
	bson.D{{"$set", bson.D{{"favorites", arrayFavorite}}}})

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK,gin.H{"message": "Remove Item from Favorite successfully"})
}

func removeIndex(s []primitive.ObjectID, index int) []primitive.ObjectID {
	return append(s[:index], s[index+1:]...)
}
