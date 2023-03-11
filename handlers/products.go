package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/service"
	validatorHelper "github.com/oleksiivelychko/go-utils/validator_helper"
	"net/http"
	"strconv"
)

// ProductHandler for CRUD actions regarding api.Product objects.
type ProductHandler struct {
	logger         hclog.Logger
	validation     *validatorHelper.Validation
	productService *service.ProductService
}

// KeyProduct is a key used for the api.Product object in the context.
type KeyProduct struct{}

func NewProductHandler(l hclog.Logger, v *validatorHelper.Validation, ps *service.ProductService) *ProductHandler {
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

type NotFound struct {
	Message string `json:"message"`
}

// GrpcError means that request to gRPC service failed.
type GrpcError struct {
	Message string `json:"message"`
}

// getProductID returns ID parameter from URL.
func (handler *ProductHandler) getProductID(r *http.Request) int {
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
