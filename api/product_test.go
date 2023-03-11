package api

import (
	"bytes"
	"fmt"
	"github.com/oleksiivelychko/go-microservice/utils"
	io "github.com/oleksiivelychko/go-utils/json_io"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductName(t *testing.T) {
	validation := utils.NewValidation()
	err := validation.Validate(Product{
		Price: 1.2,
	})

	if err == nil {
		t.Fatal("product.Name validation failed")
	}

	fmt.Println(err.Errors())
}

func TestProductPrice(t *testing.T) {
	validation := utils.NewValidation()
	err := validation.Validate(Product{
		Name:  "abc",
		Price: -1,
	})

	if err == nil {
		t.Fatal("product.Price validation failed")
	}

	fmt.Println(err.Errors())
}

func TestProductSKU(t *testing.T) {
	validation := utils.NewValidation()
	err := validation.Validate(Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "123-456-789",
	})

	if err == nil {
		t.Fatal("product.SKU validation failed")
	}

	fmt.Println(err.Errors())
}

func TestValidProduct(t *testing.T) {
	validation := utils.NewValidation()
	err := validation.Validate(Product{
		Name:  "abc",
		Price: 1.2,
		SKU:   "123-456-789",
	})

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
	err := io.ToJSON(productList, b)
	assert.NoError(t, err)
}
