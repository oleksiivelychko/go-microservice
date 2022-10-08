package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route POST /products products createProduct
// Create a new product.
//
// responses:
// 201: productResponse
// 422: validationErrorsResponse
func (p *ProductHandler) CreateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Printf("[DEBUG] POST `/products`")

	// fetch the product from the context
	product := r.Context().Value(KeyProduct{}).(*api.Product)

	api.AddProduct(product)
	rw.WriteHeader(http.StatusCreated)

	_ = utils.ToJSON(product, rw)
}
