package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route PUT /products/{id} products updateProduct
// Update a product details.
//
// responses:
// 200: productResponse
// 404: errorResponse
// 422: validationErrorsResponse
func (p *ProductHandler) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(*api.Product)
	product.ID = getProductID(r)
	p.l.Printf("[DEBUG] PUT `/products/%d`", product.ID)

	err := api.UpdateProduct(product)
	if err == api.ErrProductNotFound {
		p.l.Printf("[ERROR] PUT `/products/%d` got '%s'", product.ID, err)
		rw.WriteHeader(http.StatusNotFound)

		_ = utils.ToJSON(&NotFound{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
	_ = utils.ToJSON(product, rw)
}
