// Product API
// Documentation
//
// Schemes: http
// BasePath: /
// vVersion: 1.0.0
//
// Consumes:
// - application/json
// Produces:
// - application/json
// swagger:meta
package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"regexp"
	"time"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedAt   string  `json:"-"`
	UpdatedAt   string  `json:"-"`
	DeletedAt   string  `json:"-"`
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`[0-9]+-[0-9]+-[0-9]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}

	return true
}

type Products []*Product

/*ToJSON as JSON serializer
this approach provides better performance (reduces allocations) than json.Unmarshall()
it does not have to buffer the output into an in memory slice of bytes
*/
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func getNextID() int {
	lp := productsList[len(productsList)-1]
	return lp.ID + 1
}

func findProduct(id int) (*Product, int, error) {
	for pos, p := range productsList {
		if p.ID == id {
			return p, pos, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

func AddProducts(p *Product) {
	p.ID = getNextID()
	productsList = append(productsList, p)
}

func UpdateProducts(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	productsList[pos] = p

	return nil
}

func GetProducts() Products {
	return productsList
}

var ErrProductNotFound = fmt.Errorf("product not found")

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
