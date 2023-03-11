package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	io "github.com/oleksiivelychko/go-utils/json_io"
	"net/http"
)

// swagger:route POST /products products createProduct
// Create a new product.
//
// responses:
// 201: productResponse
// 400: grpcResponseWrapper
// 422: validationErrorsResponse
func (handler *ProductHandler) CreateProduct(writer http.ResponseWriter, request *http.Request) {
	handler.logger.Debug("POST /products CreateProduct")

	// fetch the product from the context
	product := request.Context().Value(KeyProduct{}).(*api.Product)

	err := handler.productService.AddProduct(product)
	if err != nil {
		handler.logger.Error("request to gRPC service", "error", err)
		writer.WriteHeader(http.StatusBadRequest)
		_ = io.ToJSON(&GrpcError{Message: err.Error()}, writer)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	io.ToJSON(product, writer)
}
