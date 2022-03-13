package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"net/http"
)

// swagger:route POST /products
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse
func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(api.Product)
	p.L.Printf("[DEBUG] Inserting product: %#v\n", product)

	api.AddProduct(product)
}
