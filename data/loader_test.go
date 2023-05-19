package data

import (
	"testing"
)

func TestData_LoadProductsFromJSON(t *testing.T) {
	products := LoadProductsFromJSON("./../public/data/products.json")
	if len(products) == 0 {
		t.Error("products list is empty")
	}
}
