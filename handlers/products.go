package handlers

import (
	"github.com/gorilla/mux"
	gService "github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
	"github.com/oleksiivelychko/go-microservice/utils"
	"log"
	"net/http"
	"strconv"
)

// ProductHandler for getting and updating products
type ProductHandler struct {
	l  *log.Logger
	v  *utils.Validation
	cc gService.CurrencyClient
}

// KeyProduct is a key used for the Product object in the context
type KeyProduct struct{}

// NewProductHandler returns a new product handler injected with logger, validation and gRPC client
func NewProductHandler(l *log.Logger, v *utils.Validation, cc gService.CurrencyClient) *ProductHandler {
	return &ProductHandler{l, v, cc}
}

// GenericError is a generic error message returned by a server
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

// getProductID returns the product ID from the URL
func getProductID(r *http.Request) int {
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
