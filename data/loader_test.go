package data

import (
	"testing"
)

func TestData_LoadProductsFromJSON(t *testing.T) {
	products := LoadProductsFromJson("./products.json")
	if len(products) == 0 {
		t.Error("unable to fetch, list is empty")
	}
}
