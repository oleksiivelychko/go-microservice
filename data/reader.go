package data

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/api"
	fileReader "github.com/oleksiivelychko/go-utils/json_file_reader"
)

func LoadProductsFromJson(jsonFilename string) []*api.Product {
	bytes, _ := fileReader.ReadJsonFile(jsonFilename)
	var items []*api.Product

	err := json.Unmarshal(bytes, &items)
	if err != nil {
		panic(err)
	}

	return items
}
