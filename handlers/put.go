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
func (handler *ProductHandler) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Debug("PUT /products UpdateProduct")

	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(*api.Product)
	product.ID = handler.getProductID(r)

	err := handler.productService.UpdateProduct(product)

	switch e := err.(type) {
	case *utils.GrpcServiceErr:
		handler.logger.Error("request to gRPC service", "error", err)
		rw.WriteHeader(http.StatusBadRequest)
		_ = utils.ToJSON(&GrpcError{Message: err.Error()}, rw)
		return
	case *utils.ProductNotFoundErr:
		handler.logger.Error("product not found", "id", product.ID)
		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&NotFound{Message: e.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	utils.ToJSON(product, rw)
}
