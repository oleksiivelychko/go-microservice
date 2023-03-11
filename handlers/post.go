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
func (handler *ProductHandler) CreateProduct(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Debug("POST /products CreateProduct")

	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(*api.Product)

	err := handler.productService.AddProduct(product)
	if err != nil {
		handler.logger.Error("request to gRPC service", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		_ = utils.ToJSON(&GrpcError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	utils.ToJSON(product, rw)
}
