package api

import (
	"bytes"
	"fmt"
	"github.com/oleksiivelychko/go-microservice/utils"
	"github.com/oleksiivelychko/go-utils/serializer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPI_ValidateProductName(t *testing.T) {
	validation := utils.NewValidation()
	err := validation.Validate(Product{
		Price: 1.2,
	})

	if err == nil {
		t.Fatal("unable to validate product.Name")
	}

	t.Log(err.Errors())
}

func TestAPI_ValidateProductPrice(t *testing.T) {
	validation := utils.NewValidation()
	err := validation.Validate(Product{
		Name:  "abc",
		Price: -1,
	})

	if err == nil {
		t.Fatal("unable to validate product.Price")
	}

	t.Log(err.Errors())
}

func TestAPI_ValidateProductSKU(t *testing.T) {
	validation := utils.NewValidation()
	err := validation.Validate(Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "123-456-789",
	})

	if err == nil {
		t.Fatal("unable to validate product.SKU")
	}

	fmt.Println(err.Errors())
}

func TestAPI_ValidateProduct(t *testing.T) {
	validation := utils.NewValidation()
	err := validation.Validate(Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "123-456-789",
	})

	if len(err) > 0 {
		t.Error(err.Errors())
	}
}

func TestAPI_ProductsToJSON(t *testing.T) {
	products := []*Product{
		{
			Name: "abc",
		},
	}

	bufferString := bytes.NewBufferString("")
	err := serializer.ToJSON(products, bufferString)
	assert.NoError(t, err)
}
