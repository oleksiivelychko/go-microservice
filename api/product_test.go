package api

import (
	"bytes"
	"github.com/oleksiivelychko/go-microservice/utils/serializer"
	"github.com/oleksiivelychko/go-microservice/utils/validation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAPI_ValidateProductName(t *testing.T) {
	validate := validation.New()
	err := validate.Validate(Product{
		Price: 1.2,
	})

	if err == nil {
		t.Fatal("unable to validate product.Name")
	}

	t.Log(err.Errors())
}

func TestAPI_ValidateProductPrice(t *testing.T) {
	validate := validation.New()
	err := validate.Validate(Product{
		Name:  "abc",
		Price: -1,
	})

	if err == nil {
		t.Fatal("unable to validate product.Price")
	}

	t.Log(err.Errors())
}

func TestAPI_ValidateProductSKU(t *testing.T) {
	validate := validation.New()
	err := validate.Validate(Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "123-456-789",
	})

	if err == nil {
		t.Fatal("unable to validate product.SKU")
	}

	t.Log(err.Errors())
}

func TestAPI_ValidateProduct(t *testing.T) {
	validate := validation.New()
	err := validate.Validate(Product{
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

	bufStr := bytes.NewBufferString("")
	err := serializer.ToJSON(products, bufStr)
	assert.NoError(t, err)
}
