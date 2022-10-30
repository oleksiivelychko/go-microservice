package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
	"strconv"
)

// ProductHandler for getting and updating products.
type ProductHandler struct {
	l  hclog.Logger
	v  *utils.Validation
	ps *service.ProductService
}

// KeyProduct is a key used for the Product object in the context.
type KeyProduct struct{}

// NewProductHandler returns a new product handler injected with logger, validation and gRPC client.
func NewProductHandler(l hclog.Logger, v *utils.Validation, ps *service.ProductService) *ProductHandler {
	return &ProductHandler{l, v, ps}
}

// GenericError is a generic error message returned by a server.
type GenericError struct {
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation error messages.
type ValidationErrors struct {
	Messages []string `json:"messages"`
}

// NotFound means that record not found inside collection.
type NotFound struct {
	Message string `json:"message"`
}

// GrpcError means that request to gRPC service is failed.
type GrpcError struct {
	Message string `json:"message"`
}

// getProductID returns the product ID from the URL.
func (ph *ProductHandler) getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)
	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen as the router ensures that this is a valid number
		panic(err)
	}

	return id
}

// setCurrency returns the product ID from the URL.
func (ph *ProductHandler) setCurrency(r *http.Request) {
	currency := r.URL.Query().Get("currency")
	if currency != "" {
		ph.ps.Currency.SetCurrency(currency)
	}
}
