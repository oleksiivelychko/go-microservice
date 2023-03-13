package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	jsonUtils "github.com/oleksiivelychko/go-utils/json_io"
	"net/http"
)

// swagger:route POST /products products createProduct
// Create a new product.
//
// responses:
// 201: productResponse
// 400: grpcResponseWrapper
// 422: validationErrorsResponse
func (productHandler *ProductHandler) CreateProduct(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Debug("POST /products CreateProduct")

	// fetch the product from the context
	product := request.Context().Value(KeyProduct{}).(*api.Product)

	err := productHandler.productService.AddProduct(product)
	if err != nil {
		productHandler.logger.Error("request to gRPC service", "error", err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		_ = jsonUtils.ToJSON(&GrpcError{Message: err.Error()}, responseWriter)
		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
	jsonUtils.ToJSON(product, responseWriter)
}
