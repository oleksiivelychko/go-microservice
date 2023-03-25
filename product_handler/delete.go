package product_handler

import (
	"github.com/oleksiivelychko/go-utils/io_json"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product.
//
// responses:
// 204: noContentResponse
// 404: notFoundResponse
func (productHandler *ProductHandler) DeleteProduct(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Debug("DELETE /products")
	responseWriter.Header().Add("Content-Type", "application/json")

	id := productHandler.getProductID(request)

	productNotFoundErr := productHandler.productService.DeleteProduct(id)
	if productNotFoundErr != nil {
		productHandler.logger.Error(productNotFoundErr.Error())
		responseWriter.WriteHeader(http.StatusNotFound)
		io_json.ToJSON(&productNotFoundErr, responseWriter)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}
