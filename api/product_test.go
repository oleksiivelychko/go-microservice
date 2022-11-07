package api

import (
	"bytes"
	"fmt"
	"github.com/oleksiivelychko/go-microservice/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductName(t *testing.T) {
	p := Product{
		Price: 1.2,
	}

	v := utils.NewValidation()
	err := v.Validate(p)

	if err == nil {
		t.Fatal("product.Name validation failed")
	}

	fmt.Println(err.Errors())
}

func TestProductPrice(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: -1,
	}

	v := utils.NewValidation()
	err := v.Validate(p)

	if err == nil {
		t.Fatal("product.Price validation failed")
	}

	fmt.Println(err.Errors())
}

func TestProductSKU(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "123-456-789",
	}

	v := utils.NewValidation()
	err := v.Validate(p)

	if err == nil {
		t.Fatal("product.SKU validation failed")
	}

	fmt.Println(err.Errors())
}

func TestValidProduct(t *testing.T) {
	p := Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "123-456-789",
	}

	v := utils.NewValidation()
	err := v.Validate(p)

	if len(err) > 0 {
		t.Fatal(err.Errors())
	}
}

func TestProductsToJSON(t *testing.T) {
	productList := []*Product{
		{
			Name: "abc",
		},
	}

	b := bytes.NewBufferString("")
	err := utils.ToJSON(productList, b)
	assert.NoError(t, err)
}
