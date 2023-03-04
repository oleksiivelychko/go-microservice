package data

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/api"
	"io"
	"os"
	"path/filepath"
)

func LoadDataFromJson(jsonFilename string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Open(filepath.Join(wd, jsonFilename))
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	return io.ReadAll(jsonFile)
}

func LoadProductsFromJson(jsonFilename string) []*api.Product {
	bytes, _ := LoadDataFromJson(jsonFilename)
	var items []*api.Product
	err := json.Unmarshal(bytes, &items)
	if err != nil {
		panic(err)
	}

	return items
}
