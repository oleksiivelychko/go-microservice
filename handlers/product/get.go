package product

import (
	"github.com/oleksiivelychko/go-microservice/errors"
	"github.com/oleksiivelychko/go-microservice/utils/header"
	"github.com/oleksiivelychko/go-microservice/utils/serializer"
	"net/http"
)

// swagger:route GET /products products getProducts
// Returns a list of products.
//
// responses:
// 200: productsResponse
// 400: grpcErrorResponse
func (handler *Handler) GetAll(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Debug("LIST /products")
	header.ContentTypeJSON(resp)

	products, grpcServiceErr := handler.productService.GetProducts()
	if grpcServiceErr != nil {
		handler.logger.Error(grpcServiceErr.Error())
		resp.WriteHeader(http.StatusBadRequest)
		serializer.ToJSON(&grpcServiceErr, resp)
		return
	}

	if serializerErr := serializer.ToJSON(products, resp); serializerErr != nil {
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
	handler.logger.Debug("GET /products")
	header.ContentTypeJSON(resp)

	id := handler.getProductID(req)
	product, err := handler.productService.GetProduct(id)

	switch errType := err.(type) {
	case *errors.GRPCServiceError:
		handler.logger.Error(errType.Error())
		resp.WriteHeader(http.StatusBadRequest)
		serializer.ToJSON(&errType, resp)
		return
	case *errors.ProductNotFoundError:
		handler.logger.Error(errType.Error())
		resp.WriteHeader(http.StatusNotFound)
		serializer.ToJSON(&errType, resp)
		return
	}

	if serializerErr := serializer.ToJSON(product, resp); serializerErr != nil {
		handler.logger.Error("serializer", "error", serializerErr)
	}
}
