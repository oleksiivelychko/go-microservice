package handlers

import (
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route GET /products products getProducts
// Returns a list of products.
//
// responses:
// 200: productsResponse
// 400: grpcResponseWrapper
func (handler *ProductHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Debug("GET /products GetAll")
	rw.Header().Add("Content-Type", "application/json")

	list, err := handler.productService.GetProducts()
	if err != nil {
		handler.logger.Error("request to gRPC service", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		_ = utils.ToJSON(&GrpcError{Message: err.Error()}, rw)
		return
	}

	if err = utils.ToJSON(list, rw); err != nil {
		handler.logger.Error("JSON encode", "error", err)
	}
}

// swagger:route GET /products/{id} products getProduct
// Returns a single product.
//
// responses:
// 200: productResponse
// 400: grpcResponseWrapper
// 404: notFoundResponse
// 500: errorResponse
func (handler *ProductHandler) GetOne(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Debug("GET /products GetOne")
	rw.Header().Add("Content-Type", "application/json")

	id := handler.getProductID(r)
	product, err := handler.productService.GetProduct(id)

	switch e := err.(type) {
	case *utils.GrpcServiceErr:
		handler.logger.Error("request to gRPC service", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		_ = utils.ToJSON(&GrpcError{Message: e.Error()}, rw)
		return
	case *utils.ProductNotFoundErr:
		handler.logger.Error("product not found", "id", id)
		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&NotFound{Message: e.Error()}, rw)
		return
	}

	if err = utils.ToJSON(product, rw); err != nil {
		handler.logger.Error("JSON encode", "error", err)
	}
}
