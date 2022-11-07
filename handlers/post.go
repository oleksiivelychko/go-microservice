package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route POST /products products createProduct
// Create a new product.
//
// responses:
// 201: productResponse
// 400: grpcResponseWrapper
// 422: validationErrorsResponse
func (ph *ProductHandler) CreateProduct(rw http.ResponseWriter, r *http.Request) {
	ph.log.Debug("POST /products CreateProduct")

	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(*api.Product)

	err := ph.srv.AddProduct(product)
	if err != nil {
		ph.log.Error("request to gRPC service", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		_ = utils.ToJSON(&GrpcError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	utils.ToJSON(product, rw)
}
