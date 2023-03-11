package handlers

import (
	"github.com/oleksiivelychko/go-microservice/utils"
	io "github.com/oleksiivelychko/go-utils/json_io"
	"net/http"
)

// swagger:route GET /products products getProducts
// Returns a list of products.
//
// responses:
// 200: productsResponse
// 400: grpcResponseWrapper
func (handler *ProductHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	handler.logger.Debug("GET /products GetAll")
	writer.Header().Add("Content-Type", "application/json")

	list, err := handler.productService.GetProducts()
	if err != nil {
		handler.logger.Error("request to gRPC service", "error", err)
		writer.WriteHeader(http.StatusBadRequest)
		_ = io.ToJSON(&GrpcError{Message: err.Error()}, writer)
		return
	}

	if err = io.ToJSON(list, writer); err != nil {
		handler.logger.Error("JSON encode", "error", err)
	}
}

// swagger:route GET /products/{id} products getProduct
// Returns a single product.
//
// responses:
// 200: productResponse
// 400: grpcResponseWrapper
// 404: notFoundResponse
// 500: errorResponse
func (handler *ProductHandler) GetOne(writer http.ResponseWriter, request *http.Request) {
	handler.logger.Debug("GET /products GetOne")
	writer.Header().Add("Content-Type", "application/json")

	id := handler.getProductID(request)
	product, err := handler.productService.GetProduct(id)

	switch e := err.(type) {
	case *utils.GrpcServiceErr:
		handler.logger.Error("request to gRPC service", "error", err)
		writer.WriteHeader(http.StatusBadRequest)
		_ = io.ToJSON(&GrpcError{Message: e.Error()}, writer)
		return
	case *utils.ProductNotFoundErr:
		handler.logger.Error("product not found", "id", id)
		writer.WriteHeader(http.StatusNotFound)
		_ = io.ToJSON(&NotFound{Message: e.Error()}, writer)
		return
	}

	if err = io.ToJSON(product, writer); err != nil {
		handler.logger.Error("JSON encode", "error", err)
	}
}
