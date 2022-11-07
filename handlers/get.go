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
	rw.Header().Add("Content-Type", "application/json")
	ph.l.Debug("GET GetAll /products")

	list, err := ph.ps.GetProducts()
	if err != nil {
		ph.l.Error("grpc_service.Currency.MakeExchange", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		_ = utils.ToJSON(&GrpcError{Message: err.Error()}, rw)
		return
	}

	err = utils.ToJSON(list, rw)
	if err != nil {
		ph.l.Error("serialization", "error", err)
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
	rw.Header().Add("Content-Type", "application/json")

	id := ph.getProductID(r)
	ph.l.Debug("GET GetOne /products")

	product, err := ph.ps.GetProduct(id)

	switch e := err.(type) {
	case *utils.GrpcServiceErr:
		ph.l.Error("grpc_service.Currency.MakeExchange", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		_ = utils.ToJSON(&GrpcError{Message: e.Error()}, rw)
		return
	case *utils.ProductNotFoundErr:
		ph.l.Error("product not found", "id", id)
		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&NotFound{Message: e.Error()}, rw)
		return
	}

	err = utils.ToJSON(product, rw)
	if err != nil {
		ph.l.Error("serialization", "error", err)
	}
}
