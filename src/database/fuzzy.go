package database

import (
	"fmt"
	"log"
	"strings"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	utf8string "golang.org/x/exp/utf8string"
)

// based on 
// https://medium.com/xeneta/fuzzy-search-with-mongodb-and-python-57103928ee5d
// converts string into ngrams

func CreateNgram (rawphrase string, min_size int) string{

	ngrams_map := make(map[string]bool)
	var ngrams string

	// generate for each word

	words := strings.Split(rawphrase, " ")
	for _, rawWord := range words {

		word   := utf8string.NewString(rawWord)
		length := word.RuneCount()

		for i := 0; i <= length - min_size; i++{
			for j := min_size; j <= length - i; j++{
				ngrams_map[ word.Slice(i, i + j) ] = false
			}
		}
	}
	
	for key, _ := range ngrams_map {
		ngrams += " " + key
    }
	
	//fmt.Printf("%s\n", ngrams)
	return ngrams
}

// save an Ngram for a single product

func SaveProductNgram(product Product){

	var mc *mongo.Client = GetClient()
	coll := mc.Database("buenavida").Collection("products-search")
	
	// create ngram field
	filter := bson.D{{"_id", product.Id}}
	update := bson.D{{"$set", bson.D{{"ngram", CreateNgram(product.Title, 3)}}}}

	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	
	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return
	}
}

// save an Ngram for every product in the collection

func PopulateNgrams(){

	var mc *mongo.Client = GetClient()
	coll := mc.Database("buenavida").Collection("products-search")
	
	// query all
	cursor, err := coll.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		log.Fatal("Could not execute query")
		return
	}

	// get array
	var products []Product
	if err = cursor.All(context.TODO(), &products); err != nil {
		log.Fatal("Could not execute query")
		return
	}

	for _, product := range products {
		SaveProductNgram(product)
	}
}

// Create index for Ngram field

func CreateIndexNgram(){
	var mc *mongo.Client = GetClient()
	coll := mc.Database("buenavida").Collection("products-search")
	indexView := coll.Indexes()

	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{"ngram", "text"}},
			Options: options.Index().SetName("NgramIndex"),
		},
	}

	// MaxTime limit the amount of time
	// the operation can run on the server

	opts := options.CreateIndexes().SetMaxTime(60 * time.Second)
	names, err := indexView.CreateMany(context.TODO(), models, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created indexes %v\n", names)
}
