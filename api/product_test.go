package api

import (
	"bytes"
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPI_ValidateProductName(t *testing.T) {
	validate, err := validation.New()
	if err != nil {
		t.Fatal(err)
	}

	validatorErrors := validate.Validate(Product{
		Price: 1.2,
	})

	if err == nil {
		t.Fatal("unable to validate product.Name")
	}

	t.Log(validatorErrors.Errors())
}

func TestAPI_ValidateProductPrice(t *testing.T) {
	validate, err := validation.New()
	if err != nil {
		t.Fatal(err)
	}

	validatorErrors := validate.Validate(Product{
		Name:  "abc",
		Price: -1,
	})

	if validatorErrors == nil {
		t.Fatal("unable to validate product.Price")
	}

	t.Log(validatorErrors.Errors())
}

func TestAPI_ValidateProductSKU(t *testing.T) {
	validate, err := validation.New()
	if err != nil {
		t.Fatal(err)
	}

	validatorErrors := validate.Validate(Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "123-456-789",
	})

	if validatorErrors == nil {
		t.Fatal("unable to validate product.SKU")
	}

	t.Log(validatorErrors.Errors())
}

func TestAPI_ValidateProduct(t *testing.T) {
	validate, err := validation.New()
	if err != nil {
		t.Fatal(err)
	}

	validatorErrors := validate.Validate(Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "123-456-789",
	})

	if len(validatorErrors) > 0 {
		t.Error(validatorErrors.Errors())
	}
}

func TestAPI_ProductsToJSON(t *testing.T) {
	products := []*Product{
		{
			Name: "abc",
		},
	}

	bufStr := bytes.NewBufferString("")
	err := json.NewEncoder(bufStr).Encode(products)
	assert.NoError(t, err)
}
