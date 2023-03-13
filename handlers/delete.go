package handlers

import (
	jsonUtils "github.com/oleksiivelychko/go-utils/json_io"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product.
//
// responses:
// 204: noContentResponse
// 404: notFoundResponse
func (productHandler *ProductHandler) DeleteProduct(responseWriter http.ResponseWriter, request *http.Request) {
	productHandler.logger.Debug("DELETE /products DeleteProduct")
	responseWriter.Header().Add("Content-Type", "application/json")

	id := productHandler.getProductID(request)

	err := productHandler.productService.DeleteProduct(id)
	if err != nil {
		productHandler.logger.Error("product not found", "id", id)
		responseWriter.WriteHeader(http.StatusNotFound)
		_ = jsonUtils.ToJSON(&NotFound{Message: err.Error()}, responseWriter)
		return
	}

	responseWriter.WriteHeader(http.StatusNoContent)
}
