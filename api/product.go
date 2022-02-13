package api

import (
	"encoding/json"
	"fmt"
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
