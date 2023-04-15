package product

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/errors"
	"github.com/oleksiivelychko/go-microservice/utils/header"
	"github.com/oleksiivelychko/go-microservice/utils/serializer"
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
	handler.logger.Debug("PUT /products")
	header.ContentTypeJSON(resp)

	product := req.Context().Value(KeyProduct{}).(*api.Product)
	product.ID = handler.getProductID(req)

	err := handler.productService.UpdateProduct(product)

	switch errType := err.(type) {
	case *errors.GRPCServiceError:
		handler.logger.Error(errType.Error())
		resp.WriteHeader(http.StatusBadRequest)
		serializer.ToJSON(&errType, resp)
		return
	case *errors.ProductNotFoundError:
		handler.logger.Error(errType.Error())
		resp.WriteHeader(http.StatusNotFound)
		serializer.ToJSON(&errType, resp)
		return
	}

	resp.WriteHeader(http.StatusOK)
	serializer.ToJSON(product, resp)
}
