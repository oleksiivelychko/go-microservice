package product

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/errors"
	"net/http"
)

// swagger:route GET /products products getProducts
// Returns a list of products.
//
// responses:
// 200: productsResponse
// 400: grpcErrorResponse
func (handler *Handler) GetAll(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Info("LIST /products")
	resp.Header().Set("Content-Type", "application/json")

	products, grpcServiceErr := handler.productService.GetProducts()
	if grpcServiceErr != nil {
		handler.logger.Error(grpcServiceErr.Error())
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(&grpcServiceErr)
		return
	}

	if serializerErr := json.NewEncoder(resp).Encode(products); serializerErr != nil {
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
func (handler *Handler) GetOne(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Info("GET /products")
	resp.Header().Set("Content-Type", "application/json")

	id := handler.getProductID(req)
	product, err := handler.productService.GetProduct(id)

	switch errType := err.(type) {
	case *errors.GRPCServiceError:
		handler.logger.Error(errType.Error())
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(&errType)
		return
	case *errors.ProductNotFoundError:
		handler.logger.Error(errType.Error())
		resp.WriteHeader(http.StatusNotFound)
		json.NewEncoder(resp).Encode(&errType)
		return
	}

	if serializerErr := json.NewEncoder(resp).Encode(product); serializerErr != nil {
		handler.logger.Error("serializer", "error", serializerErr)
	}
}
