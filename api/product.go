package api

import (
	"github.com/oleksiivelychko/go-microservice/utils/datetime"
)

// Product structure for an API model.
// swagger:model product
type Product struct {
	// required: false
	// min: 1
	ID int `json:"id"`

	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"gt=0"`

	// SKU - in the field of inventory management, a stock keeping unit is a distinct type of item for sale, purchased, or tracked in inventory,
	// such as a product or service, and all attributes associated with the item type that distinguish it from other item types.
	// For a product, these attributes can include manufacturer, description, material, size, color, packaging, and warranty terms.
	// When a business takes inventory of its stock, it counts the quantity it has of each SKU.
	// SKU can also refer to a unique identifier or code, sometimes represented via a barcode for scanning and tracking, that refers to the particular stock keeping unit.
	// These identifiers are not regulated or standardized.
	// When a company receives items from a vendor, it has a choice of maintaining the vendor's SKU or creating its own.
	//
	// Original source: https://en.wikipedia.org/wiki/Stock_keeping_unit
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`

	// required: false
	UpdatedAt datetime.JSON `json:"updatedAt,omitempty"`
}

type Products []*Product
