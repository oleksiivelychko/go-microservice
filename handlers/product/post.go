package product

import (
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/api"
	"net/http"
)

// swagger:route POST /products products createProduct
// Create a new product.
//
// responses:
// 201: productResponse
// 400: grpcErrorResponse
// 422: validationErrorsResponse
func (handler *Handler) CreateProduct(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Info("POST /products")
	resp.Header().Set("Content-Type", "application/json")

	product := req.Context().Value(KeyProduct{}).(*api.Product)

	grpcServiceErr := handler.productService.AddProduct(product)
	if grpcServiceErr != nil {
		handler.logger.Error(grpcServiceErr.Error())
		resp.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(resp).Encode(&grpcServiceErr)
		return
	}

	resp.WriteHeader(http.StatusCreated)
	json.NewEncoder(resp).Encode(product)
}
