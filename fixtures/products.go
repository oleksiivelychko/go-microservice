package fixtures

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"time"
)

var productsList = []*api.Product{
	&api.Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       1.49,
		SKU:         "LATTE-01",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
	&api.Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       0.99,
		SKU:         "ESPRESSO-01",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
}
