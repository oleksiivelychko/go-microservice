package product

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils/header"
	"github.com/oleksiivelychko/go-microservice/utils/serializer"
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
	handler.logger.Debug("POST /products")
	header.ContentTypeJSON(resp)

	product := req.Context().Value(KeyProduct{}).(*api.Product)

	grpcServiceErr := handler.productService.AddProduct(product)
	if grpcServiceErr != nil {
		handler.logger.Error(grpcServiceErr.Error())
		resp.WriteHeader(http.StatusBadRequest)
		serializer.ToJSON(&grpcServiceErr, resp)
		return
	}

	resp.WriteHeader(http.StatusCreated)
	serializer.ToJSON(product, resp)
}
