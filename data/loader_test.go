package data

import (
	"testing"
)

func TestLoadProductsFromJSON(t *testing.T) {
	products := LoadProductsFromJson("./products.json")
	if len(products) == 0 {
		t.Error("products list is empty")
	}
}
