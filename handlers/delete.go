package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product.
//
// responses:
// 204: noContentResponse
// 404: notFoundResponse
func (p *ProductHandler) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)
	p.l.Printf("[DEBUG] DELETE `/products/%d`", id)

	err := api.DeleteProduct(id)
	if err == api.ErrProductNotFound {
		p.l.Printf("[ERROR] DELETE `/products/%d` got '%s'", id, err)
		rw.WriteHeader(http.StatusNotFound)

		_ = utils.ToJSON(&NotFound{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
