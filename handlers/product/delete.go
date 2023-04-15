package product

import (
	"github.com/oleksiivelychko/go-microservice/utils/header"
	"github.com/oleksiivelychko/go-microservice/utils/serializer"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product.
//
// responses:
// 204: noContentResponse
// 404: notFoundResponse
func (handler *Handler) DeleteProduct(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Debug("DELETE /products")
	header.ContentTypeJSON(resp)

	id := handler.getProductID(req)

	productNotFoundErr := handler.productService.DeleteProduct(id)
	if productNotFoundErr != nil {
		handler.logger.Error(productNotFoundErr.Error())
		resp.WriteHeader(http.StatusNotFound)
		serializer.ToJSON(&productNotFoundErr, resp)
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}
