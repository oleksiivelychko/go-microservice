package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	io "github.com/oleksiivelychko/go-utils/json_io"
	"net/http"
)

// swagger:route PUT /products/{id} products updateProduct
// Update a product details.
//
// responses:
// 200: productResponse
// 400: grpcResponseWrapper
// 404: errorResponse
// 422: validationErrorsResponse
func (handler *ProductHandler) UpdateProduct(writer http.ResponseWriter, request *http.Request) {
	handler.logger.Debug("PUT /products UpdateProduct")

	// fetch the product from the context
	product := request.Context().Value(KeyProduct{}).(*api.Product)
	product.ID = handler.getProductID(request)

	err := handler.productService.UpdateProduct(product)

	switch e := err.(type) {
	case *utils.GrpcServiceErr:
		handler.logger.Error("request to gRPC service", "error", err)
		writer.WriteHeader(http.StatusBadRequest)
		_ = io.ToJSON(&GrpcError{Message: err.Error()}, writer)
		return
	case *utils.ProductNotFoundErr:
		handler.logger.Error("product not found", "id", product.ID)
		writer.WriteHeader(http.StatusNotFound)
		_ = io.ToJSON(&NotFound{Message: e.Error()}, writer)
		return
	}

	writer.WriteHeader(http.StatusOK)
	io.ToJSON(product, writer)
}
