package data

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/api"
	"io"
	"os"
	"path/filepath"
)

func LoadProductsFromJSON() []*api.Product {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Open(filepath.Join(wd, "./public/products.json"))
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var products api.Products
	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		panic(err)
	}

	return products
}
