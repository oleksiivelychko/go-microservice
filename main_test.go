package main

import (
	"fmt"
	"github.com/oleksiivelychko/go-helper/env_addr"
	"github.com/oleksiivelychko/go-helper/pretty_bytes"
	httpClient "github.com/oleksiivelychko/go-microservice/sdk/client"
	"github.com/oleksiivelychko/go-microservice/sdk/client/products"
	"testing"
)

/**
main server must be running before
*/

func TestHttpClientGetProducts(t *testing.T) {
	client := createHttpClient()

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
	client := createHttpClient()

	params := products.NewGetProductParams()
	params.ID = 1
	productOne, err := client.Products.GetProduct(params)

	if err != nil {
		t.Fatal(err)
	}

	p, _ := productOne.GetPayload().MarshalBinary()
	out := pretty_bytes.PrettyBytes(p, "	")
	fmt.Printf("%s\n", out)

}

func createHttpClient() *httpClient.GoMicroservice {
	addr := env_addr.GetAddr()
	cfg := httpClient.DefaultTransportConfig().WithHost(addr)
	return httpClient.NewHTTPClientWithConfig(nil, cfg)
}
