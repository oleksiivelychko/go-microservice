package main

import (
	"fmt"
	"github.com/oleksiivelychko/go-helper/env_addr"
	"github.com/oleksiivelychko/go-helper/pretty_bytes"
	httpClient "github.com/oleksiivelychko/go-microservice/sdk/client"
	"github.com/oleksiivelychko/go-microservice/sdk/client/products"
	"github.com/oleksiivelychko/go-microservice/sdk/models"
	"testing"
)

/**
main server must be running before
*/

var client = createHttpClient()

func TestHttpClientGetProducts(t *testing.T) {
	params := products.NewGetProductsParams()
	productsList, err := client.Products.GetProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	for _, productItem := range productsList.GetPayload() {
		p, _ := productItem.MarshalBinary()
		out := pretty_bytes.PrettyBytes(p, "	")
		fmt.Printf("%s\n", out)
	}
}

func TestHttpClientGetProduct(t *testing.T) {
	productOne, err := fetchProduct(1)

	if err != nil {
		t.Fatal(err)
	}

	p, _ := productOne.GetPayload().MarshalBinary()
	out := pretty_bytes.PrettyBytes(p, "	")
	fmt.Printf("%s\n", out)
}

func TestHttpClientCreateProduct(t *testing.T) {
	params := products.NewCreateProductParams()

	var pName = "Coffee"
	var pPrice float32 = 1.49
	var pSKU = "000-000-000"
	var pDescription = "Coffee with milk"
	params.Body = &models.Product{
		ID:          3,
		Name:        &pName,
		Description: pDescription,
		Price:       &pPrice,
		SKU:         &pSKU,
	}

	_, err := client.Products.CreateProduct(params)
	if err != nil {
		t.Fatal(err)
	}

	productOne, err := fetchProduct(3)
	if err != nil {
		t.Fatal(err)
	}

	if *productOne.GetPayload().Name != pName {
		t.Fatal("Product name doesn't math")
	}

	if productOne.GetPayload().Description != pDescription {
		t.Fatal("Product description doesn't math")
	}

	if *productOne.GetPayload().Price != pPrice {
		t.Fatal("Product price doesn't math")
	}

	if *productOne.GetPayload().SKU != pSKU {
		t.Fatal("Product SKU doesn't math")
	}
}

/**
TestHttpClientUpdateProduct
https://github.com/go-swagger/go-swagger/discussions/2742
*/
func TestHttpClientUpdateProduct(t *testing.T) {

}

func TestHttpClientDeleteProduct(t *testing.T) {

}

func createHttpClient() *httpClient.GoMicroservice {
	addr := env_addr.GetAddr()
	cfg := httpClient.DefaultTransportConfig().WithHost(addr)
	return httpClient.NewHTTPClientWithConfig(nil, cfg)
}

func fetchProduct(id int64) (*products.GetProductOK, error) {
	params := products.NewGetProductParams()
	params.ID = id
	return client.Products.GetProduct(params)
}
