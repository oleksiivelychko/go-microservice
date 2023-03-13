package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	jsonUtils "github.com/oleksiivelychko/go-utils/json_io"
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
func (productHandler *ProductHandler) UpdateProduct(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Debug("PUT /products UpdateProduct")

	// fetch the product from the context
	product := request.Context().Value(KeyProduct{}).(*api.Product)
	product.ID = productHandler.getProductID(request)

	err := productHandler.productService.UpdateProduct(product)

	switch e := err.(type) {
	case *utils.GrpcServiceErr:
		productHandler.logger.Error("request to gRPC service", "error", err)
		responseWriter.WriteHeader(http.StatusBadRequest)
		_ = jsonUtils.ToJSON(&GrpcError{Message: err.Error()}, responseWriter)
		return
	case *utils.ProductNotFoundErr:
		productHandler.logger.Error("product not found", "id", product.ID)
		responseWriter.WriteHeader(http.StatusNotFound)
		_ = jsonUtils.ToJSON(&NotFound{Message: e.Error()}, responseWriter)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	jsonUtils.ToJSON(product, responseWriter)
}
