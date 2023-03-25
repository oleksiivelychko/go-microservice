package product_handler

import (
	"fmt"
	"github.com/oleksiivelychko/go-microservice/errors"
	"github.com/oleksiivelychko/go-microservice/utils"
	"github.com/oleksiivelychko/go-utils/io_json"
	"net/http"
)

// swagger:route GET /products products getProducts
// Returns a list of products.
//
// responses:
// 200: productsResponse
// 400: grpcErrorResponse
func (productHandler *ProductHandler) GetAll(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Debug(fmt.Sprintf("LIST %s", utils.ProductsURL))
	responseWriter.Header().Add("Content-Type", "application/json")

	products, grpcServiceErr := productHandler.productService.GetProducts()
	if grpcServiceErr != nil {
		productHandler.logger.Error(grpcServiceErr.Error())
		responseWriter.WriteHeader(http.StatusBadRequest)
		_ = io_json.ToJSON(&grpcServiceErr, responseWriter)
		return
	}

	if serializerErr := io_json.ToJSON(products, responseWriter); serializerErr != nil {
		productHandler.logger.Error("serializer", "error", serializerErr)
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
func (productHandler *ProductHandler) GetOne(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Debug(fmt.Sprintf("GET %s", utils.ProductsURL))
	responseWriter.Header().Add("Content-Type", "application/json")

	id := productHandler.getProductID(request)
	product, err := productHandler.productService.GetProduct(id)

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

	if serializerErr := io_json.ToJSON(product, responseWriter); serializerErr != nil {
		productHandler.logger.Error("serializer", "error", serializerErr)
	}
}
