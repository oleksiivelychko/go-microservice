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
func (ph *ProductHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	ph.log.Debug("GET /products GetAll")
	rw.Header().Add("Content-Type", "application/json")

	list, err := ph.srv.GetProducts()
	if err != nil {
		ph.log.Error("request to gRPC service", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		_ = utils.ToJSON(&GrpcError{Message: err.Error()}, rw)
		return
	}

	if err = utils.ToJSON(list, rw); err != nil {
		ph.log.Error("serialization", "error", err)
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
func (ph *ProductHandler) GetOne(rw http.ResponseWriter, r *http.Request) {
	ph.log.Debug("GET /products GetOne")
	rw.Header().Add("Content-Type", "application/json")

	id := ph.getProductID(r)
	product, err := ph.srv.GetProduct(id)

	switch e := err.(type) {
	case *utils.GrpcServiceErr:
		ph.log.Error("request to gRPC service", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		_ = utils.ToJSON(&GrpcError{Message: e.Error()}, rw)
		return
	case *utils.ProductNotFoundErr:
		ph.log.Error("product not found", "id", id)
		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&NotFound{Message: e.Error()}, rw)
		return
	}

	if err = utils.ToJSON(product, rw); err != nil {
		ph.log.Error("serialization", "error", err)
	}
}
