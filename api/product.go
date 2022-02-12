package api

import (
	"encoding/json"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedAt   string  `json:"-"`
	UpdatedAt   string  `json:"-"`
	DeletedAt   string  `json:"-"`
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return []*Product{
		&Product{
			ID:          1,
			Name:        "Latte",
			Description: "Frothy milky coffee",
			Price:       1.49,
			SKU:         "LATTE-01",
			CreatedAt:   time.Now().UTC().String(),
			UpdatedAt:   time.Now().UTC().String(),
		},
		&Product{
			ID:          2,
			Name:        "Espresso",
			Description: "Short and strong coffee without milk",
			Price:       0.99,
			SKU:         "ESPRESSO-01",
			CreatedAt:   time.Now().UTC().String(),
			UpdatedAt:   time.Now().UTC().String(),
		},
	}
}
