package main

import (
	"fmt"
	"github.com/oleksiivelychko/go-microservice/sdk/client"
	"github.com/oleksiivelychko/go-microservice/sdk/client/products"
	"github.com/oleksiivelychko/go-microservice/sdk/models"
	"github.com/oleksiivelychko/go-microservice/utils"
	"github.com/oleksiivelychko/go-utils/indent_json"
	"testing"
)

/*
Warning: main and gRPC servers must be running before.
*/
var sdkClient = makeClient()

func TestMain_GetProducts(t *testing.T) {
	productsParams := products.NewGetProductsParams()
	productsList, err := sdkClient.Products.GetProducts(productsParams)

	if err != nil {
		t.Fatal(err)
	}

	for _, productItem := range productsList.GetPayload() {
		p, _ := productItem.MarshalBinary()
		out := indent_json.IndentJSON(string(p))
		fmt.Printf("%s\n", out)
	}
}

func TestMain_GetProduct(t *testing.T) {
	product, err := fetchProduct(1)

	if err != nil {
		t.Fatal(err)
	}

	productBytes, _ := product.GetPayload().MarshalBinary()
	out := indent_json.IndentJSON(string(productBytes))
	fmt.Printf("%s\n", out)
}

func TestMain_CreateProduct(t *testing.T) {
	productParams := products.NewCreateProductParams()

	var productName = "Coffee"
	var productPrice = 1.49
	var productSKU = "000-000-000"

	productParams.Body = &models.Product{
		Name:  &productName,
		Price: &productPrice,
		SKU:   &productSKU,
	}

	_, err := sdkClient.Products.CreateProduct(productParams)
	if err != nil {
		t.Fatal(err)
	}

	product, err := fetchProduct(3)
	if err != nil {
		t.Fatal(err)
	}

	if *product.GetPayload().Name != productName {
		t.Errorf("product.Name '%s' from payload doesn't match the '%s'", *product.GetPayload().Name, productName)
	}

	if *product.GetPayload().SKU != productSKU {
		t.Errorf("product.SKU '%s' from payload doesn't match the '%s'", *product.GetPayload().SKU, productSKU)
	}

	if *product.GetPayload().Price > productPrice {
		t.Errorf("product.Price '%f' from payload is greater than '%f'", *product.GetPayload().Price, productPrice)
	}
}

/*
TestMain_UpdateProduct
https://github.com/go-swagger/go-swagger/discussions/2742
*/
func TestMain_UpdateProduct(t *testing.T) {
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

	_, err := sdkClient.Products.UpdateProduct(productParams)
	if err != nil {
		t.Fatal(err)
	}

	product, err := fetchProduct(3)
	if err != nil {
		t.Fatal(err)
	}

	if *product.GetPayload().Name != productName {
		t.Errorf("product.Name '%s' from payload doesn't match the '%s'", *product.GetPayload().Name, productName)
	}

	if *product.GetPayload().SKU != productSKU {
		t.Errorf("product.SKU '%s' from payload doesn't match the '%s'", *product.GetPayload().SKU, productSKU)
	}

	if *product.GetPayload().Price > productPrice {
		t.Errorf("product.Price '%f' from payload is greater than '%f'", *product.GetPayload().Price, productPrice)
	}
}

func TestMain_DeleteProduct(t *testing.T) {
	productParams := products.NewDeleteProductParams()
	productParams.ID = 3

	_, err := sdkClient.Products.DeleteProduct(productParams)
	if err != nil {
		t.Fatal(err)
	}

	product, err := fetchProduct(3)
	if product != nil {
		t.Fatal(err)
	}
}

func makeClient() *client.GoMicroservice {
	return client.NewHTTPClientWithConfig(nil, client.DefaultTransportConfig().WithHost(utils.ServerAddr))
}

func fetchProduct(id int64) (*products.GetProductOK, error) {
	productParams := products.NewGetProductParams()
	productParams.ID = id
	return sdkClient.Products.GetProduct(productParams)
}
