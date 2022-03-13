package api

import (
	"fmt"
	"time"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// Unique identifier
	//
	// required: false
	// min: 1
	ID int `json:"id"`

	// name
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// description
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// price
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"gt=0"`

	// SKU - in the field of inventory management, a stock keeping unit is a distinct type of item for sale, purchased, or tracked in inventory,
	// such as a product or service, and all attributes associated with the item type that distinguish it from other item types.
	// For a product, these attributes can include manufacturer, description, material, size, color, packaging, and warranty terms.
	// When a business takes inventory of its stock, it counts the quantity it has of each SKU.
	// SKU can also refer to a unique identifier or code, sometimes represented via a barcode for scanning and tracking, that refers to the particular stock keeping unit.
	// These identifiers are not regulated or standardized.
	// When a company receives items from a vendor, it has a choice of maintaining the vendor's SKU or creating its own
	// Original source: https://en.wikipedia.org/wiki/Stock_keeping_unit
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`

	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

type Products []*Product

var ErrProductNotFound = fmt.Errorf("product not found")

func GetProducts() Products {
	return productsList
}

func GetProduct(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	return productsList[i], nil
}

func UpdateProduct(p Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	productsList[i] = &p
	return nil
}

func AddProduct(p Product) {
	// get the next id in sequence
	maxID := productsList[len(productsList)-1].ID
	p.ID = maxID + 1
	productsList = append(productsList, &p)
}

func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productsList = append(productsList[:i], productsList[i+1])
	return nil
}

/**
findIndexByProductID finds the index of a product
returns -1 when no product can be found
*/
func findIndexByProductID(id int) int {
	for i, p := range productsList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

var productsList = []*Product{
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
