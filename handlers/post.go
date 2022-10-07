package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"net/http"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  400: errorValidation
//  422: errorValidation
//  501: errorResponse
func (p *ProductHandler) CreateProduct(rw http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(*api.Product)
	p.l.Printf("[DEBUG] create a new product: %#v\n", product)

	api.AddProduct(*product)
}
