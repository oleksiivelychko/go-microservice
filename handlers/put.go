package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route PUT /products/{id} products updateProduct
// Update a product details.
//
// responses:
// 200: productResponse
// 400: grpcResponseWrapper
// 404: errorResponse
// 422: validationErrorsResponse
func (ph *ProductHandler) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	ph.log.Debug("PUT UpdateProduct /products")
	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(*api.Product)
	product.ID = ph.getProductID(r)

	err := ph.srv.UpdateProduct(product)

	switch e := err.(type) {
	case *utils.GrpcServiceErr:
		ph.log.Error("grpc_service.Currency.MakeExchange", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		_ = utils.ToJSON(&GrpcError{Message: err.Error()}, rw)
		return
	case *utils.ProductNotFoundErr:
		ph.log.Error("product not found", "id", product.ID)
		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&NotFound{Message: e.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_ = utils.ToJSON(product, rw)
}
