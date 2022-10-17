package handlers

import (
	"context"
	gService "github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route GET /products products getProducts
// Returns a list of products.
//
// responses:
// 200: productsResponse
func (p *ProductHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	p.l.Printf("[DEBUG] GET `/products`")

	list := api.GetProducts()

	err := utils.ToJSON(list, rw)
	if err != nil {
		p.l.Printf("[ERROR] GET `/products` during serialization got '%s'", err)
	}
}

// swagger:route GET /products/{id} products getProduct
// Returns a product by ID.
//
// responses:
// 200: productResponse
// 404: notFoundResponse
// 500: errorResponse
func (p *ProductHandler) GetOne(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)
	p.l.Printf("[DEBUG] GET `/products/%d`", id)

	product, err := api.GetProduct(id)

	switch err {
	case nil:
	case api.ErrProductNotFound:
		p.l.Printf("[ERROR] GET `/products/%d` got '%s'", id, err)

		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&NotFound{Message: err.Error()}, rw)
		return
	default:
		p.l.Printf("[ERROR] GET `/products/%d` got internal server error '%s'", id, err)

		rw.WriteHeader(http.StatusInternalServerError)
		_ = utils.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	er := &gService.ExchangeRequest{
		From: gService.Currencies_USD.String(),
		To:   gService.Currencies_EUR.String(),
	}
	rateResponse, err := p.cc.MakeExchange(context.Background(), er)
	if err == nil {
		p.l.Printf("[INFO] GET `grpc_service.Currency.MakeExchange` got rate=%f", rateResponse.Rate)
		product.Price *= rateResponse.Rate
	} else {
		p.l.Printf("[ERROR] GET `grpc_service.Currency.MakeExchange` got '%s'", err)
	}

	err = utils.ToJSON(product, rw)
	if err != nil {
		p.l.Printf("[ERROR] GET `/products/%d` during serialization got '%s'", id, err)
	}
}
