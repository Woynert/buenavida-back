package database

// TEMPORALY DATABASE 
// THIS SHOULD BE REPLACED WITH MONGODB AND POSTGRESQL
// ONLY FOR TEST PORPUSES

type Product struct {
	Id				string	`json:"id"`
	Title			string	`json:"title"`
	Units			string	`json:"units"`
	Price			float64	`json:"price"`
	Discount		int		`json:"discount"`
	PricePerUnit	float64	`json:"priceperunit"`
	Description		string	`json:"description"`

}

var products = []Product{
	{
      Id: "1",
      Title: "Parches de Oro de 24 kt Rejuvenecedores para Contorno de Ojos",
      Units: "60UDS.",
      Price: 15.50,
	  Discount: 50,
      PricePerUnit: 0.26,
      Description: "Parches de oro de 24 kt",
    },
    {
      Id: "2",
      Title: "Parches Iluminadores para el Contorno de Ojos",
      Units: "60UDS.",
      Price: 15.50,
	  Discount: 50,
      PricePerUnit: 0.26,
      Description: "Parches iluminadores",
    },
    {
      Id: "3",
      Title: "Parches Supertonificantes para Contorno de Ojos",
      Units: "60UDS.",
      Price: 15.50,
	  Discount: 50,
      PricePerUnit: 0.26,
      Description: "Parches supertonificantes",
    },
}

func GetFilteredProducts(searchterm, minprice, maxprice, pageid string) ([]Product, error){
	return products, nil
}

func CartGetItems() ([]Product, error){
	return products, nil
}
