package data

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-utils/file_reader_json"
)

func LoadProductsFromJson(filename string) []*api.Product {
	bytes, _ := file_reader_json.FileReaderJSON(filename)
	var items []*api.Product

	err := json.Unmarshal(bytes, &items)
	if err != nil {
		panic(err)
	}

	return items
}
