package main

import (
	"github.com/oleksiivelychko/go-microservice/sdk/client"
	"github.com/oleksiivelychko/go-microservice/sdk/client/products"
	"github.com/oleksiivelychko/go-microservice/sdk/models"
	"github.com/oleksiivelychko/go-utils/formatter"
	"testing"
)

/*
Warning: main and gRPC servers must be running before.
*/
var sdkClient = makeClient("localhost:9090")

func TestMain_GetProducts(t *testing.T) {
	productsList, err := sdkClient.Products.GetProducts(products.NewGetProductsParams())
	if err != nil {
		t.Fatal(err)
	}

	for _, productItem := range productsList.GetPayload() {
		productBytes, _ := productItem.MarshalBinary()
		productJSON := formatter.IndentJSON(string(productBytes))
		t.Logf("%s\n", productJSON)
	}
}

func TestMain_GetProduct(t *testing.T) {
	product, err := fetchProduct(1)
	if err != nil {
		t.Fatal(err)
	}

	productBytes, _ := product.GetPayload().MarshalBinary()
	productJSON := formatter.IndentJSON(string(productBytes))
	t.Logf("%s\n", productJSON)
}

func TestMain_CreateProduct(t *testing.T) {
	params := products.NewCreateProductParams()

	var name = "Coffee"
	var price = 1.49
	var SKU = "000-000-000"

	params.Body = &models.Product{
		Name:  &name,
		Price: &price,
		SKU:   &SKU,
	}

	_, err := sdkClient.Products.CreateProduct(params)
	if err != nil {
		t.Fatal(err)
	}

	product, err := fetchProduct(3)
	if err != nil {
		t.Fatal(err)
	}

	if *product.GetPayload().Name != name {
		t.Errorf("product.Name %s from payload does not match the %s", *product.GetPayload().Name, name)
	}

	if *product.GetPayload().SKU != SKU {
		t.Errorf("product.SKU %s from payload does not match the %s", *product.GetPayload().SKU, SKU)
	}

	if *product.GetPayload().Price == price {
		t.Errorf("product.Price %f from payload equals to %f", *product.GetPayload().Price, price)
	}
}

/*
TestMain_UpdateProduct
https://github.com/go-swagger/go-swagger/discussions/2742
*/
func TestMain_UpdateProduct(t *testing.T) {
	params := products.NewUpdateProductParams()

	var name = "Coffee with milk"
	var price = 1.99
	var SKU = "111-111-111"

	params.ID = 3
	params.Body = &models.Product{
		Name:  &name,
		Price: &price,
		SKU:   &SKU,
	}

	_, err := sdkClient.Products.UpdateProduct(params)
	if err != nil {
		t.Fatal(err)
	}

	product, err := fetchProduct(3)
	if err != nil {
		t.Fatal(err)
	}

	if *product.GetPayload().Name != name {
		t.Errorf("product.Name %s from payload does not match the %s", *product.GetPayload().Name, name)
	}

	if *product.GetPayload().SKU != SKU {
		t.Errorf("product.SKU '%s' from payload does not match the %s", *product.GetPayload().SKU, SKU)
	}

	if *product.GetPayload().Price == price {
		t.Errorf("product.Price %f from payload equals to %f", *product.GetPayload().Price, price)
	}
}

func TestMain_DeleteProduct(t *testing.T) {
	params := products.NewDeleteProductParams()
	params.ID = 3

	_, err := sdkClient.Products.DeleteProduct(params)
	if err != nil {
		t.Fatal(err)
	}

	product, err := fetchProduct(3)
	if product != nil {
		t.Fatal(err)
	}
}

func makeClient(addr string) *client.GoMicroservice {
	return client.NewHTTPClientWithConfig(nil, client.DefaultTransportConfig().WithHost(addr))
}

func fetchProduct(id int64) (*products.GetProductOK, error) {
	params := products.NewGetProductParams()
	params.ID = id
	return sdkClient.Products.GetProduct(params)
}
