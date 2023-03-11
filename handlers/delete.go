package handlers

import (
	io "github.com/oleksiivelychko/go-utils/json_io"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product.
//
// responses:
// 204: noContentResponse
// 404: notFoundResponse
func (handler *ProductHandler) DeleteProduct(writer http.ResponseWriter, request *http.Request) {
	handler.logger.Debug("DELETE /products DeleteProduct")
	writer.Header().Add("Content-Type", "application/json")

	id := handler.getProductID(request)

	err := handler.productService.DeleteProduct(id)
	if err != nil {
		handler.logger.Error("product not found", "id", id)
		writer.WriteHeader(http.StatusNotFound)
		_ = io.ToJSON(&NotFound{Message: err.Error()}, writer)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
