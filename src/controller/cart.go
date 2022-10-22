package controller

import (
	db "woynert/buenavida-api/database"

	//"fmt"
    "net/http"
    "github.com/gin-gonic/gin"
)

func CartGetItems(c *gin.Context) {
	var products []db.Product
	products, err := db.CartGetItems()

	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError,
			gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, products)
}
