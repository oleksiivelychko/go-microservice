package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route PUT /products
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation
func (p *ProductHandler) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(api.Product)
	p.L.Println("[DEBUG] updating record id", product.ID)

	err := api.UpdateProduct(product)
	if err == api.ErrProductNotFound {
		p.L.Println("[ERROR] product not found", err)

		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&GenericError{Message: "Product not found"}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
