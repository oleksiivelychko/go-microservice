package handler

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/errors"
	"github.com/oleksiivelychko/go-utils/response"
	"github.com/oleksiivelychko/go-utils/serializer"
	"net/http"
)

// swagger:route PUT /products/{id} products updateProduct
// Update a product details.
//
// responses:
// 200: productResponse
// 400: grpcErrorResponse
// 404: errorResponse
// 422: validationErrorsResponse
func (handler *ProductHandler) UpdateProduct(responseWriter http.ResponseWriter, request *http.Request) {
	response.HeaderContentTypeJSON(responseWriter)
	handler.logger.Debug("PUT /products")

	product := request.Context().Value(KeyProduct{}).(*api.Product)
	product.ID = handler.getProductID(request)

	err := handler.productService.UpdateProduct(product)

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

	responseWriter.WriteHeader(http.StatusOK)
	serializer.ToJSON(product, responseWriter)
}
