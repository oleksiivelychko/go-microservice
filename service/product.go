package service

import (
	"fmt"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/data"
	"github.com/oleksiivelychko/go-microservice/utils"
)

type ProductService struct {
	currency *CurrencyService
	data     []*api.Product
}

func NewProductService(currency *CurrencyService) *ProductService {
	var productsList = data.LoadProductsFromJSON()
	return &ProductService{currency, productsList}
}

func (ps *ProductService) GetProducts() (api.Products, error) {
	rate, err := ps.currency.GetRate()
	if err != nil {
		return ps.data, &utils.GrpcServiceRequestErr{Err: err.Error()}
	}

	ratedProductsList := api.Products{}
	for _, product := range ps.data {
		ratedProduct := *product
		ratedProduct.Price *= rate
		ratedProductsList = append(ratedProductsList, &ratedProduct)
	}

	return ratedProductsList, nil
}

func (ps *ProductService) GetProduct(id int) (*api.Product, error) {
	i := ps.findIndexByProductID(id)
	if i == -1 {
		return nil, &utils.ProductNotFoundErr{Err: fmt.Sprintf("id=%d", id)}
	}

	rate, err := ps.currency.GetRate()
	if err != nil {
		return ps.data[i], &utils.GrpcServiceRequestErr{Err: err.Error()}
	}

	ratedProduct := *ps.data[i]
	ratedProduct.Price *= rate

	return &ratedProduct, nil
}

func (ps *ProductService) AddProduct(p *api.Product) error {
	p.ID = ps.getNextProductId()

	rate, err := ps.currency.GetRate()
	if err != nil {
		err = &utils.GrpcServiceRequestErr{Err: err.Error()}
	} else {
		p.Price *= rate
	}

	ps.data = append(ps.data, p)
	return err
}

func (ps *ProductService) UpdateProduct(p *api.Product) error {
	i := ps.findIndexByProductID(p.ID)
	if i == -1 {
		return &utils.ProductNotFoundErr{Err: fmt.Sprintf("id=%d", p.ID)}
	}

	rate, err := ps.currency.GetRate()
	if err != nil {
		err = &utils.GrpcServiceRequestErr{Err: err.Error()}
	} else {
		p.Price *= rate
	}

	ps.data[i] = p
	return err
}

func (ps *ProductService) DeleteProduct(id int) error {
	index := ps.findIndexByProductID(id)
	if index == -1 {
		return &utils.ProductNotFoundErr{Err: fmt.Sprintf("id=%d", id)}
	}

	ps.data = ps.deleteProductByIndex(index)
	return nil
}

func (ps *ProductService) getNextProductId() int {
	if len(ps.data) == 0 {
		return 1
	}

	return ps.data[len(ps.data)-1].ID + 1
}

/*
*
findIndexByProductID finds the index of a product.
Returns -1 when product not found.
*/
func (ps *ProductService) findIndexByProductID(id int) int {
	for i, p := range ps.data {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (ps *ProductService) deleteProductByIndex(index int) []*api.Product {
	return append(ps.data[:index], ps.data[index+1:]...)
}