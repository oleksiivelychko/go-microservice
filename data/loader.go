package data

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/api"
	"os"
)

func LoadProductsFromJSON(filename string) []*api.Product {
	bytesArr, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var products []*api.Product

	err = json.Unmarshal(bytesArr, &products)
	if err != nil {
		panic(err)
	}

	return products
}
