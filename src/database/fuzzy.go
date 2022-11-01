package database

import (
	"fmt"
	"log"
	"strings"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

// https://medium.com/xeneta/fuzzy-search-with-mongodb-and-python-57103928ee5d

func MakeNgram(phrase string, min_size int) string{
	ngrams_map := make(map[string]bool)
	var ngrams []string

	words := strings.Split(phrase, " ")
	
	for _, word := range words {

		length := len(word)

		for i := 0; i <= length - min_size; i++{
			for j := min_size; j <= length - i; j++{
				
				// fmt.Printf("%s\n", word[i:i + j])
				ngrams_map[ word[i:i + j] ] = false
			}
		}
	}
	
	for key, _ := range ngrams_map {
		ngrams = append(ngrams, key)
    }
	
	fmt.Printf("%s\n", ngrams)
	//return ngrams
	return strings.Join(ngrams," ")
}

func IndexNgrams(product Product){

	// save an ngram for every item in products-search

	var mc *mongo.Client = GetClient()
	coll := mc.Database("buenavida").Collection("products-search")

	
	//update := bson.D{{"$inc", bson.D{{"price", 100}}}}
	filter := bson.D{{"_id", product.Id}}
	update := bson.D{{"$set", bson.D{{"ngrams", MakeNgram(product.Title, 3)}}}}

	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	
	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return
	}
}

func UpdateAllNgram(){

	var mc *mongo.Client = GetClient()
	coll := mc.Database("buenavida").Collection("products-search")
	opts := options.Find().SetSort(bson.D{{"_id", 1}}) // sort by id
	
	// query
	cursor, err := coll.Find(context.TODO(), bson.D{{}}, opts)
	if err != nil {
		log.Fatal("Could not execute query")
		return
	}

	// Get a list of all returned documents and print them out.
	var products []Product
	if err = cursor.All(context.TODO(), &products); err != nil {
		log.Fatal("Could not execute query")
		return
	}

	for _, product := range products {
		//fmt.Println(product)
		IndexNgrams(product)
	}
}
