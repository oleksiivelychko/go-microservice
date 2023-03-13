package handlers

import (
	"github.com/oleksiivelychko/go-microservice/utils"
	jsonUtils "github.com/oleksiivelychko/go-utils/json_io"
	"net/http"
)

// swagger:route GET /products products getProducts
// Returns a list of products.
//
// responses:
// 200: productsResponse
// 400: grpcResponseWrapper
func (productHandler *ProductHandler) GetAll(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Debug("GET /products GetAll")
	responseWriter.Header().Add("Content-Type", "application/json")

	products, err := productHandler.productService.GetProducts()
	if err != nil {
		productHandler.logger.Error("request to gRPC service", "error", err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		_ = jsonUtils.ToJSON(&GrpcError{Message: err.Error()}, responseWriter)
		return
	}

	if err = jsonUtils.ToJSON(products, responseWriter); err != nil {
		productHandler.logger.Error("JSON encode", "error", err)
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
func (productHandler *ProductHandler) GetOne(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Debug("GET /products GetOne")
	responseWriter.Header().Add("Content-Type", "application/json")

	id := productHandler.getProductID(request)
	product, err := productHandler.productService.GetProduct(id)

	switch e := err.(type) {
	case *utils.GrpcServiceErr:
		productHandler.logger.Error("request to gRPC service", "error", err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		_ = jsonUtils.ToJSON(&GrpcError{Message: e.Error()}, responseWriter)
		return
	case *utils.ProductNotFoundErr:
		productHandler.logger.Error("product not found", "id", id)
		responseWriter.WriteHeader(http.StatusNotFound)
		_ = jsonUtils.ToJSON(&NotFound{Message: e.Error()}, responseWriter)
		return
	}

	if err = jsonUtils.ToJSON(product, responseWriter); err != nil {
		productHandler.logger.Error("JSON encode", "error", err)
	}
}
