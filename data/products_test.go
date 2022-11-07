package data

import (
	"testing"
)

func TestLoadProductsFromJSON(t *testing.T) {
	products := LoadProductsFromJSON("./../public/products.json")
	if len(products) == 0 {
		t.Fatal("products list is empty")
	}
}
