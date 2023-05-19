package product

import (
	"github.com/gorilla/mux"
	"github.com/oleksiivelychko/go-grpc-service/logger"
	"github.com/oleksiivelychko/go-microservice/services"
	"github.com/oleksiivelychko/go-microservice/validation"
	"net/http"
	"strconv"
)

// KeyProduct is a key used for the api.Product object in the context.
type KeyProduct struct{}

// Handler for CRUD actions regarding api.Product objects.
type Handler struct {
	logger         *logger.Logger
	validation     *validation.Validate
	productService *services.Product
}

func NewHandler(validation *validation.Validate, productService *services.Product, logger *logger.Logger) *Handler {
	return &Handler{logger, validation, productService}
}

// getProductID returns ID parameter from URL.
func (handler *Handler) getProductID(r *http.Request) int {
	muxVars := mux.Vars(r)
	id, err := strconv.Atoi(muxVars["id"])
	if err != nil {
		// should never happen as the router ensures that this is a valid number
		panic(err)
	}

	return id
}
