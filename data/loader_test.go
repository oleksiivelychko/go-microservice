package data

import (
	"testing"
)

func TestData_LoadProductsFromJSON(t *testing.T) {
	products := LoadProductsFromJSON("./../public/data/products.json")
	if len(products) == 0 {
		t.Error("unable to load, list is empty")
	}
}
