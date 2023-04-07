package handler

import (
	"fmt"
	"github.com/oleksiivelychko/go-microservice/utils"
	"github.com/oleksiivelychko/go-utils/response"
	"github.com/oleksiivelychko/go-utils/serializer"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product.
//
// responses:
// 204: noContentResponse
// 404: notFoundResponse
func (handler *ProductHandler) DeleteProduct(responseWriter http.ResponseWriter, request *http.Request) {
	handler.logger.Debug(fmt.Sprintf("DELETE %s", utils.ProductsURL))
	response.HeaderContentTypeJSON(responseWriter)

	id := handler.getProductID(request)

	productNotFoundErr := handler.productService.DeleteProduct(id)
	if productNotFoundErr != nil {
		handler.logger.Error(productNotFoundErr.Error())
		responseWriter.WriteHeader(http.StatusNotFound)
		serializer.ToJSON(&productNotFoundErr, responseWriter)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}
