package api

import (
	"bytes"
	"fmt"
	"github.com/oleksiivelychko/go-microservice/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductName(t *testing.T) {
	p := Product{
		Price: 1.2,
	}

	v := helpers.NewValidation()
	err := v.Validate(p)

	if err == nil {
		t.Fatal("product:name:validation didn't pass test")
	}

	fmt.Println(err.Errors())
}

func TestProductPrice(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: -1,
	}

	v := helpers.NewValidation()
	err := v.Validate(p)

	if err == nil {
		t.Fatal("product:price:validation didn't pass test")
	}

	fmt.Println(err.Errors())
}

func TestProductSKU(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "abc",
	}

	v := helpers.NewValidation()
	err := v.Validate(p)

	if err == nil {
		t.Fatal("product:sku:validation didn't pass test")
	}

	fmt.Println(err.Errors())
}

func TestValidProduct(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "123-456-789",
	}

	v := helpers.NewValidation()
	err := v.Validate(p)

	if len(err) > 0 {
		t.Fatal(err.Errors())
	}
}

func TestProductsToJSON(t *testing.T) {
	ps := []*Product{
		&Product{
			Name: "abc",
		},
	}

	b := bytes.NewBufferString("")
	err := helpers.ToJSON(ps, b)
	assert.NoError(t, err)
}
