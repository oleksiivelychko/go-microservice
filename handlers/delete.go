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
func (ph *ProductHandler) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	ph.l.Debug("DeleteProduct")
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)

	err := ph.ps.DeleteProduct(id)
	if err != nil {
		ph.l.Debug("product not found", "id", id)
		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&NotFound{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
