package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id				primitive.ObjectID `json:"_id"          bson:"_id"`   
	Title			string             `json:"title"        bson:"title"`
	Units			string             `json:"units"        bson:"units"`
	Price			float64            `json:"price"        bson:"price"`
	Discount		int	               `json:"discount"     bson:"discount"`
	PricePerUnit	string             `json:"priceperunit" bson:"priceperunit"`
	Description		string             `json:"description"  bson:"description"`
	Score           float64            `json:"score"        bson:"score"`
	ImageUrl        string             `json:"imageurl"     bson:"imageurl"`
}

type User struct {
	Id        primitive.ObjectID   `json:"_id"       bson:"_id"`
	Firstname string               `json:"firstname" bson:"firstname"`
	Lastname  string               `json:"lastname"  bson:"lastname"`
	Email     string               `json:"email"     bson:"email"`
	Password  string               `json:"password"  bson:"password"`
	Favorites []primitive.ObjectID `json:"favorites" bson:"favorites"`
}

