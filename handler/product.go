package handler

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-utils/validation"
	"net/http"
	"strconv"
)

// KeyProduct is a key used for the api.Product object in the context.
type KeyProduct struct{}

// ProductHandler for CRUD actions regarding api.Product objects.
type ProductHandler struct {
	logger         hclog.Logger
	validation     *validation.Validate
	productService *service.ProductService
}

func NewProductHandler(logger hclog.Logger, validation *validation.Validate, productService *service.ProductService) *ProductHandler {
	return &ProductHandler{logger, validation, productService}
}

// getProductID returns ID parameter from URL.
func (handler *ProductHandler) getProductID(r *http.Request) int {
	muxVars := mux.Vars(r)
	id, err := strconv.Atoi(muxVars["id"])
	if err != nil {
		// should never happen as the router ensures that this is a valid number
		panic(err)
	}

	return id
}
