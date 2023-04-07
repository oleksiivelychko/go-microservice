package handler

import (
	"github.com/oleksiivelychko/go-microservice/errors"
	"github.com/oleksiivelychko/go-utils/response"
	"github.com/oleksiivelychko/go-utils/serializer"
	"net/http"
)

// swagger:route GET /products products getProducts
// Returns a list of products.
//
// responses:
// 200: productsResponse
// 400: grpcErrorResponse
func (handler *ProductHandler) GetAll(responseWriter http.ResponseWriter, request *http.Request) {
	handler.logger.Debug("LIST /products")
	response.HeaderContentTypeJSON(responseWriter)

	products, grpcServiceErr := handler.productService.GetProducts()
	if grpcServiceErr != nil {
		handler.logger.Error(grpcServiceErr.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		_ = serializer.ToJSON(&grpcServiceErr, responseWriter)
		return
	}

	if serializerErr := serializer.ToJSON(products, responseWriter); serializerErr != nil {
		handler.logger.Error("serializer", "error", serializerErr)
	}
}

// swagger:route GET /products/{id} products getProduct
// Returns a single product.
//
// responses:
// 200: productResponse
// 400: grpcErrorResponse
// 404: notFoundResponse
// 500: errorResponse
func (handler *ProductHandler) GetOne(responseWriter http.ResponseWriter, request *http.Request) {
	handler.logger.Debug("GET /products")
	response.HeaderContentTypeJSON(responseWriter)

	id := handler.getProductID(request)
	product, err := handler.productService.GetProduct(id)

	switch errType := err.(type) {
	case *errors.GRPCServiceError:
		handler.logger.Error(errType.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		_ = serializer.ToJSON(&errType, responseWriter)
		return
	case *errors.ProductNotFoundError:
		handler.logger.Error(errType.Error())
		responseWriter.WriteHeader(http.StatusNotFound)
		_ = serializer.ToJSON(&errType, responseWriter)
		return
	}

	if serializerErr := serializer.ToJSON(product, responseWriter); serializerErr != nil {
		handler.logger.Error("serializer", "error", serializerErr)
	}
}
