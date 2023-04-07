package handler

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-utils/response"
	"github.com/oleksiivelychko/go-utils/serializer"
	"net/http"
)

// swagger:route POST /products products createProduct
// Create a new product.
//
// responses:
// 201: productResponse
// 400: grpcErrorResponse
// 422: validationErrorsResponse
func (handler *ProductHandler) CreateProduct(responseWriter http.ResponseWriter, request *http.Request) {
	handler.logger.Debug("POST /products")
	response.HeaderContentTypeJSON(responseWriter)

	product := request.Context().Value(KeyProduct{}).(*api.Product)

	grpcServiceErr := handler.productService.AddProduct(product)
	if grpcServiceErr != nil {
		handler.logger.Error(grpcServiceErr.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		serializer.ToJSON(&grpcServiceErr, responseWriter)
		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
	serializer.ToJSON(product, responseWriter)
}
