package controller

import (
	db "woynert/buenavida-api/database"

	"fmt"
	"context"
	"net/http"
	"strconv"
	"math"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StoreFilterItems (c *gin.Context) {

	// get url parameters
	// and sanitize parameters

	var err error
	searchTerm := c.DefaultQuery("searchterm", "")

	pageId, err := strconv.Atoi(c.DefaultQuery("pageid", "0"))
	if err != nil { pageId = 0 }
	pageId = int(math.Max(0, float64(pageId)))

	minprice, err := strconv.Atoi(c.DefaultQuery("minprice"  , "0"))
	if err != nil { minprice = 0 }
	minprice = int(math.Max(0, float64(minprice)))

	maxprice, err := strconv.Atoi(c.DefaultQuery("maxprice"  , "0"))
	if err != nil { maxprice = 0 }
	maxprice = int(math.Max(0, float64(maxprice)))

	// get connection and collection

	var mc *mongo.Client = db.MongoGetClient()
	coll := mc.Database("buenavida").Collection("products")

	opts := options.Find()
	opts = opts.SetSkip(int64(12 * pageId)) // skip previous items
	opts = opts.SetLimit(12)                // get only 12 items

	// apply filters

	filter := bson.D{
		{"price", bson.D{{ "$gte", minprice }}},
	}

	if (maxprice != 0){
		filter = append(filter, bson.E{ "price", bson.D{{ "$lte", maxprice }}})
	}

	projection := bson.D{{"ngram", 0}} // don't return ngrams

	if (searchTerm != ""){
		// sort by score
		projection = append(projection, bson.E{ "score", bson.D{{"$meta", "textScore"}} })
		opts   = opts.SetSort(bson.D{{"score", bson.D{{ "$meta", "textScore" }} }})
		filter = append(filter, bson.E{ "$text", bson.D{{ "$search", db.CreateNgram(searchTerm, 3) }} })
	} else {
		opts = opts.SetSort(bson.D{{"_id", 1}}) // sort by id
	}
	opts = opts.SetProjection(projection)

	// fmt.Println(filter)

	// query
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Could not execute query"})
		return
	}

	// get Product array
	var products []db.Product
	if err = cursor.All(context.TODO(), &products); err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Internal server error"})
		return
	}

	fmt.Printf(
	"\nFound %d products\nSearchterm \"%s\" PageId %d Minprice %d Maxprice %d\n",
	len(products), searchTerm, pageId, minprice, maxprice)

	c.IndentedJSON(http.StatusOK, products)
}
