package api

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Product defines the structure for an API model.
// swagger:model product
type Product struct {
	// unique identifier
	//
	// required: false
	// min: 1
	ID int `json:"id"`

	// name
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

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
	// When a company receives items from a vendor, it has a choice of maintaining the vendor's SKU or creating its own.
	//
	// Original source: https://en.wikipedia.org/wiki/Stock_keeping_unit
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`

	UpdatedAt DateTime `json:"updatedAt,omitempty"`
}

type DateTime struct {
	time.Time
}

func (t *DateTime) MarshalJSON() ([]byte, error) {
	stamp := time.Now().Format(time.RFC3339)
	return []byte("\"" + stamp + "\""), nil
}

func (t *DateTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	date, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.Time = date
	return
}

type Products []*Product

var ErrProductNotFound = fmt.Errorf("product not found")

var productsList = LoadProductsFromJSON()

func GetProducts() Products {
	return productsList
}

func GetProduct(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}

	return productsList[i], nil
}

func AddProduct(p *Product) {
	p.ID = GetNextProductId()
	productsList = append(productsList, p)
}

func UpdateProduct(p *Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	productsList[i] = p
	return nil
}

func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productsList = RemoveProductByIndex(productsList, i)
	return nil
}

/*
*
findIndexByProductID finds the index of a product.
Returns -1 when no product not found.
*/
func findIndexByProductID(id int) int {
	for i, p := range productsList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func RemoveProductByIndex(s []*Product, index int) []*Product {
	return append(s[:index], s[index+1:]...)
}

func GetNextProductId() int {
	if len(productsList) == 0 {
		return 1
	}

	return productsList[len(productsList)-1].ID + 1
}

func LoadProductsFromJSON() []*Product {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Open(filepath.Join(wd, "./public/products.json"))
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var products Products
	_ = json.Unmarshal(byteValue, &products)

	return products
}
