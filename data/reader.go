package data

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/api"
	jsonUtils "github.com/oleksiivelychko/go-utils/json_file_reader"
)

func LoadProductsFromJson(filename string) []*api.Product {
	bytes, _ := jsonUtils.ReadJsonFile(filename)
	var items []*api.Product

	err := json.Unmarshal(bytes, &items)
	if err != nil {
		panic(err)
	}

	return items
}
