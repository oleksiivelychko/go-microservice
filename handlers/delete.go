package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse
func (p *ProductHandler) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)
	p.L.Println("[DEBUG] deleting record id", id)

	err := api.DeleteProduct(id)
	if err == api.ErrProductNotFound {
		p.L.Println("[ERROR] deleting record id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.L.Println("[ERROR] deleting record", err)

		rw.WriteHeader(http.StatusInternalServerError)
		_ = utils.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
