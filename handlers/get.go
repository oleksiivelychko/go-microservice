package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route GET /products products listProducts
// Returns a list of products.
// responses:
//	200: productsResponse
func (p *ProductHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] get all records")

	list := api.GetProducts()
	err := utils.ToJSON(list, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
	}
}

// swagger:route GET /products/{id}
// Returns a single product by ID
// responses:
//	200: productResponse
//  404: errorResponse
func (p *ProductHandler) GetOne(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)
	p.l.Println("[DEBUG] get record id", id)

	product, err := api.GetProduct(id)

	switch err {
	case nil:
	case api.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		_ = utils.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		_ = utils.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = utils.ToJSON(product, rw)

	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
	}
}
