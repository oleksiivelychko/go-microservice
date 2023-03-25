package product_handler

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-utils/validator_helper"
	"net/http"
	"strconv"
)

// KeyProduct is a key used for the api.Product object in the context.
type KeyProduct struct{}

// ProductHandler for CRUD actions regarding api.Product objects.
type ProductHandler struct {
	logger         hclog.Logger
	validation     *validator_helper.Validation
	productService *service.ProductService
}

func NewProductHandler(
	logger hclog.Logger,
	validation *validator_helper.Validation,
	productService *service.ProductService,
) *ProductHandler {
	return &ProductHandler{logger, validation, productService}
}

// getProductID returns ID parameter from URL.
func (productHandler *ProductHandler) getProductID(r *http.Request) int {
	muxVars := mux.Vars(r)
	id, err := strconv.Atoi(muxVars["id"])
	if err != nil {
		// should never happen as the router ensures that this is a valid number
		panic(err)
	}

	return id
}
