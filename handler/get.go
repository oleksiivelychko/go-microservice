package handler

import (
	"github.com/oleksiivelychko/go-microservice/errors"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route GET /products products getProducts
// Returns a list of products.
//
// responses:
// 200: productsResponse
// 400: grpcErrorResponse
func (handler *Product) GetAll(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Debug("LIST /products")
	utils.HeaderContentTypeJSON(resp)

	products, grpcServiceErr := handler.productService.GetProducts()
	if grpcServiceErr != nil {
		handler.logger.Error(grpcServiceErr.Error())
		resp.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&grpcServiceErr, resp)
		return
	}

	if serializerErr := utils.ToJSON(products, resp); serializerErr != nil {
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
func (handler *Product) GetOne(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Debug("GET /products")
	utils.HeaderContentTypeJSON(resp)

	id := handler.getProductID(req)
	product, err := handler.productService.GetProduct(id)

	switch errType := err.(type) {
	case *errors.GRPCServiceError:
		handler.logger.Error(errType.Error())
		resp.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&errType, resp)
		return
	case *errors.ProductNotFoundError:
		handler.logger.Error(errType.Error())
		resp.WriteHeader(http.StatusNotFound)
		utils.ToJSON(&errType, resp)
		return
	}

	if serializerErr := utils.ToJSON(product, resp); serializerErr != nil {
		handler.logger.Error("serializer", "error", serializerErr)
	}
}
