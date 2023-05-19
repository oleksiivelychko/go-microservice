package product

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/errors"
	"net/http"
)

// swagger:route PUT /products/{id} products updateProduct
// Update a product details.
//
// responses:
// 200: productResponse
// 400: grpcErrorResponse
// 404: errorResponse
// 422: validationErrorsResponse
func (handler *Handler) UpdateProduct(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Info("PUT /products")
	resp.Header().Set("Content-Type", "application/json")

	respEncoder := json.NewEncoder(resp)

	product := req.Context().Value(KeyProduct{}).(*api.Product)
	product.ID = handler.getProductID(req)

	err := handler.productService.UpdateProduct(product)

	switch errType := err.(type) {
	case *errors.GRPCServiceError:
		handler.logger.Error(errType.Error())
		resp.WriteHeader(http.StatusBadRequest)
		respEncoder.Encode(&errType)
		return
	case *errors.ProductNotFoundError:
		handler.logger.Error(errType.Error())
		resp.WriteHeader(http.StatusNotFound)
		respEncoder.Encode(&errType)
		return
	}

	resp.WriteHeader(http.StatusOK)
	respEncoder.Encode(product)
}
