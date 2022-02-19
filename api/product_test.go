package api

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "Milk",
		Price: 1.00,
		SKU:   "123-456-789",
	}
	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
