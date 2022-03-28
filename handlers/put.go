package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route PUT /products/{id} products updateProduct
// Update a products details
//
// responses:
//	200: productResponse
//  404: errorResponse
//  422: errorValidation
//  501: errorResponse
func (p *ProductHandler) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(api.Product)
	p.l.Println("[DEBUG] updating record id", product.ID)

	err := api.UpdateProduct(product)
	if err == api.ErrProductNotFound {
		p.l.Println("[ERROR] updating record id does not exist", err)

		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&GenericError{Message: "product not found"}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
