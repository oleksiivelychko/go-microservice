package product_handler

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-utils/io_json"
	"net/http"
)

// swagger:route POST /products products createProduct
// Create a new product.
//
// responses:
// 201: productResponse
// 400: grpcErrorResponse
// 422: validationErrorsResponse
func (productHandler *ProductHandler) CreateProduct(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Debug("POST /products")

	product := request.Context().Value(KeyProduct{}).(*api.Product)

	grpcServiceErr := productHandler.productService.AddProduct(product)
	if grpcServiceErr != nil {
		productHandler.logger.Error(grpcServiceErr.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		io_json.ToJSON(&grpcServiceErr, responseWriter)
		return
	}

	responseWriter.WriteHeader(http.StatusCreated)
	io_json.ToJSON(product, responseWriter)
}
