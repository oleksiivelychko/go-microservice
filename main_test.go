package main

import (
	"fmt"
	httpClient "github.com/oleksiivelychko/go-microservice/sdk/client"
	"github.com/oleksiivelychko/go-microservice/sdk/client/products"
	"github.com/oleksiivelychko/go-microservice/sdk/models"
	prettily "github.com/oleksiivelychko/go-microservice/utils/echo_bytes"
	"testing"
)

/*
*
Warning: main server must be running before.
*/
var client = createHttpClient("localhost:9090")

func TestHttpClientGetProducts(t *testing.T) {
	params := products.NewGetProductsParams()
	productsList, err := client.Products.GetProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	for _, productItem := range productsList.GetPayload() {
		p, _ := productItem.MarshalBinary()
		out := prettily.EchoBytes(p, "	")
		fmt.Printf("%s\n", out)
	}
}

func TestHttpClientGetProduct(t *testing.T) {
	productOne, err := fetchProduct(1)

	if err != nil {
		t.Fatal(err)
	}

	p, _ := productOne.GetPayload().MarshalBinary()
	out := prettily.EchoBytes(p, "	")
	fmt.Printf("%s\n", out)
}

func TestHttpClientCreateProduct(t *testing.T) {
	params := products.NewCreateProductParams()

	var pName = "Coffee"
	var pPrice float32 = 1.49
	var pSKU = "000-000-000"

	params.Body = &models.Product{
		Name:  &pName,
		Price: &pPrice,
		SKU:   &pSKU,
	}

	_, err := client.Products.CreateProduct(params)
	if err != nil {
		t.Fatal(err)
	}

	product, err := fetchProduct(3)
	if err != nil {
		t.Fatal(err)
	}

	if *product.GetPayload().Name != pName {
		t.Fatal("product:name doesn't match")
	}

	if *product.GetPayload().Price != pPrice {
		t.Fatal("product:price doesn't match")
	}

	if *product.GetPayload().SKU != pSKU {
		t.Fatal("product:SKU doesn't match")
	}
}

/*
*
TestHttpClientUpdateProduct
https://github.com/go-swagger/go-swagger/discussions/2742
*/
func TestHttpClientUpdateProduct(t *testing.T) {
	params := products.NewUpdateProductParams()

	var pName = "Coffee with milk"
	var pPrice float32 = 1.99
	var pSKU = "111-111-111"

	params.ID = 3
	params.Body = &models.Product{
		Name:  &pName,
		Price: &pPrice,
		SKU:   &pSKU,
	}

	_, err := client.Products.UpdateProduct(params)
	if err != nil {
		t.Fatal(err)
	}

	product, err := fetchProduct(3)
	if err != nil {
		t.Fatal(err)
	}

	if *product.GetPayload().Name != pName {
		t.Fatal("product:name doesn't match")
	}

	if *product.GetPayload().Price != pPrice {
		t.Fatal("product:price doesn't match")
	}

	if *product.GetPayload().SKU != pSKU {
		t.Fatal("product:SKU doesn't match")
	}
}

func TestHttpClientDeleteProduct(t *testing.T) {
	params := products.NewDeleteProductParams()
	params.ID = 3

	_, err := client.Products.DeleteProduct(params)
	if err != nil {
		t.Fatal(err)
	}

	product, err := fetchProduct(3)
	if product != nil {
		t.Fatal(err)
	}
}

func createHttpClient(addr string) *httpClient.GoMicroservice {
	cfg := httpClient.DefaultTransportConfig().WithHost(addr)
	return httpClient.NewHTTPClientWithConfig(nil, cfg)
}

func fetchProduct(id int64) (*products.GetProductOK, error) {
	params := products.NewGetProductParams()
	params.ID = id
	return client.Products.GetProduct(params)
}
