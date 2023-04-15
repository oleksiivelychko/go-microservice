package handler

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-microservice/validation"
	"net/http"
	"strconv"
)

// KeyProduct is a key used for the api.Product object in the context.
type KeyProduct struct{}

// Product for CRUD actions regarding api.Product objects.
type Product struct {
	logger         hclog.Logger
	validation     *validation.Validate
	productService *service.ProductService
}

func New(logger hclog.Logger, validation *validation.Validate, productService *service.ProductService) *Product {
	return &Product{logger, validation, productService}
}

// getProductID returns ID parameter from URL.
func (handler *Product) getProductID(r *http.Request) int {
	muxVars := mux.Vars(r)
	id, err := strconv.Atoi(muxVars["id"])
	if err != nil {
		// should never happen as the router ensures that this is a valid number
		panic(err)
	}

	return id
}
