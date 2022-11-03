package controller

import (
	db "woynert/buenavida-api/database"

	"fmt"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type PaymentFormProduct struct {
	Id         primitive.ObjectID `json:"itemid"   bson:"itemid"`   
	Quantity   int     `json:"quantity"   bson:"quantity"`	
	Discount   int     `json:"discount"   bson:"discount"`	
	PriceBase  float64 `json:"pricebase"  bson:"pricebase"`	
}

func Payment(c *gin.Context) {

	// get user id

	userIdAny, exists := c.Get("userid")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Invalid token"})
		return
	}
	userId := userIdAny.(string)

	// get products from body

	var items []PaymentFormProduct

	if err := c.BindJSON(&items); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	// get discount and price for every item

	mc   := db.MongoGetClient()
	coll := mc.Database("buenavida").Collection("products")

	for index, _ := range items {

		item := &items[index]
		var product db.Product

		// query
		err := coll.FindOne(
			context.TODO(),
			bson.D{{"_id", item.Id}},
		).Decode(&product)

		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal server error"})
			return
		}

		item.Discount   = product.Discount
		item.PriceBase  = product.Price

	}

	// create sale payment 
	// get sale id

	pg := db.PgGetClient()

	var saleId string

	row := pg.QueryRow("SELECT * FROM sales_new($1)", userId)
	if err := row.Scan(&saleId); err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
		gin.H{"message": "Internal server error"})
		return
	}

	// add every item to the sale payment

	for _, item := range items {

		_, err := pg.Exec("CALL sales_add_item ($1, $2, $3, $4, $5)",
		saleId, // sale id
		item.Id.Hex(), // product id
		item.Quantity, // quantity
		item.Discount, // discount
		item.PriceBase) // price

		if err != nil {
			fmt.Println("Error inserting data ", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"message": "Internal server error"})
			return 
		}
	}

	c.IndentedJSON(http.StatusOK,
	gin.H{"message": "Payment successful"})
}
