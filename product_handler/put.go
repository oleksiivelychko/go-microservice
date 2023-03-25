package product_handler

import (
	"fmt"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/errors"
	"github.com/oleksiivelychko/go-microservice/utils"
	"github.com/oleksiivelychko/go-utils/io_json"
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
func (productHandler *ProductHandler) UpdateProduct(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Debug(fmt.Sprintf("PUT %s", utils.ProductsURL))

	product := request.Context().Value(KeyProduct{}).(*api.Product)
	product.ID = productHandler.getProductID(request)

	err := productHandler.productService.UpdateProduct(product)

	switch errType := err.(type) {
	case *errors.GRPCServiceError:
		productHandler.logger.Error(errType.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		_ = io_json.ToJSON(&errType, responseWriter)
		return
	case *errors.ProductNotFoundError:
		productHandler.logger.Error(errType.Error())
		responseWriter.WriteHeader(http.StatusNotFound)
		_ = io_json.ToJSON(&errType, responseWriter)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
	io_json.ToJSON(product, responseWriter)
}
