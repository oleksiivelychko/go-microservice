package handler

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
// 400: grpcErrorResponse
// 422: validationErrorsResponse
func (handler *Product) CreateProduct(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Debug("POST /products")
	utils.HeaderContentTypeJSON(resp)

	product := req.Context().Value(KeyProduct{}).(*api.Product)

	grpcServiceErr := handler.productService.AddProduct(product)
	if grpcServiceErr != nil {
		handler.logger.Error(grpcServiceErr.Error())
		resp.WriteHeader(http.StatusBadRequest)
		utils.ToJSON(&grpcServiceErr, resp)
		return
	}

	resp.WriteHeader(http.StatusCreated)
	utils.ToJSON(product, resp)
}
