package handlers

import (
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product.
//
// responses:
// 204: noContentResponse
// 404: notFoundResponse
func (handler *ProductHandler) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	handler.logger.Debug("DELETE /products DeleteProduct")
	rw.Header().Add("Content-Type", "application/json")

	id := handler.getProductID(r)

	err := handler.productService.DeleteProduct(id)
	if err != nil {
		handler.logger.Error("product not found", "id", id)
		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&NotFound{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
