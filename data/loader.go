package data

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils/reader"
)

func LoadProductsFromJson(filename string) []*api.Product {
	bytesArr, _ := reader.ReadFile(filename)
	var products []*api.Product

	err := json.Unmarshal(bytesArr, &products)
	if err != nil {
		panic(err)
	}

	return products
}
