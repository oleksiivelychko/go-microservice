package main

import (
	"fmt"
	httpClient "github.com/oleksiivelychko/go-microservice/sdk/client"
	"github.com/oleksiivelychko/go-microservice/sdk/client/products"
	"github.com/oleksiivelychko/go-microservice/sdk/models"
	jsonUtils "github.com/oleksiivelychko/go-utils/json_indent"
	"testing"
)

/*
*
Warning: main and gRPC servers must be running before.
*/
var client = createHttpClient("localhost:9090")

func TestHttpClientGetProducts(t *testing.T) {
	productsParams := products.NewGetProductsParams()
	productsList, err := client.Products.GetProducts(productsParams)

	if err != nil {
		t.Error(err)
	}

	for _, productItem := range productsList.GetPayload() {
		p, _ := productItem.MarshalBinary()
		out := jsonUtils.JsonIndent(string(p))
		fmt.Printf("%s\n", out)
	}
}

func TestHttpClientGetProduct(t *testing.T) {
	product, err := fetchProduct(1)

	if err != nil {
		t.Error(err)
	}

	productBytes, _ := product.GetPayload().MarshalBinary()
	out := jsonUtils.JsonIndent(string(productBytes))
	fmt.Printf("%s\n", out)
}

func TestHttpClientCreateProduct(t *testing.T) {
	productParams := products.NewCreateProductParams()

	var productName = "Coffee"
	var productPrice = 1.49
	var productSKU = "000-000-000"

	productParams.Body = &models.Product{
		Name:  &productName,
		Price: &productPrice,
		SKU:   &productSKU,
	}

	_, err := client.Products.CreateProduct(productParams)
	if err != nil {
		t.Error(err)
	}

	product, err := fetchProduct(3)
	if err != nil {
		t.Error(err)
	}

	if *product.GetPayload().Name != productName {
		t.Errorf("product.Name `%s` doesn't match `%s`", *product.GetPayload().Name, productName)
	}

	if *product.GetPayload().SKU != productSKU {
		t.Errorf("product.SKU `%s` doesn't match `%s`", *product.GetPayload().SKU, productSKU)
	}

	if *product.GetPayload().Price < productPrice {
		t.Errorf("product.Price `%f` didn't update `%f`", *product.GetPayload().Price, productPrice)
	}
}

/*
*
TestHttpClientUpdateProduct
https://github.com/go-swagger/go-swagger/discussions/2742
*/
func TestHttpClientUpdateProduct(t *testing.T) {
	productParams := products.NewUpdateProductParams()

	var productName = "Coffee with milk"
	var productPrice = 1.99
	var productSKU = "111-111-111"

	productParams.ID = 3
	productParams.Body = &models.Product{
		Name:  &productName,
		Price: &productPrice,
		SKU:   &productSKU,
	}

	_, err := client.Products.UpdateProduct(productParams)
	if err != nil {
		t.Error(err)
	}

	product, err := fetchProduct(3)
	if err != nil {
		t.Error(err)
	}

	if *product.GetPayload().Name != productName {
		t.Errorf("product.Name `%s` doesn't match `%s`", *product.GetPayload().Name, productName)
	}

	if *product.GetPayload().SKU != productSKU {
		t.Errorf("product.SKU `%s` doesn't match `%s`", *product.GetPayload().SKU, productSKU)
	}

	if *product.GetPayload().Price < productPrice {
		t.Errorf("product.Price `%f` didn't update `%f`", *product.GetPayload().Price, productPrice)
	}
}

func TestHttpClientDeleteProduct(t *testing.T) {
	productParams := products.NewDeleteProductParams()
	productParams.ID = 3

	_, err := client.Products.DeleteProduct(productParams)
	if err != nil {
		t.Error(err)
	}

	product, err := fetchProduct(3)
	if product != nil {
		t.Error(err)
	}
}

func createHttpClient(addr string) *httpClient.GoMicroservice {
	config := httpClient.DefaultTransportConfig().WithHost(addr)
	return httpClient.NewHTTPClientWithConfig(nil, config)
}

func fetchProduct(id int64) (*products.GetProductOK, error) {
	productParams := products.NewGetProductParams()
	productParams.ID = id
	return client.Products.GetProduct(productParams)
}
