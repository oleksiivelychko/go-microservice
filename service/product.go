package service

import (
	"fmt"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/data"
	"github.com/oleksiivelychko/go-microservice/utils"
)

type ProductService struct {
	CurrencyService *CurrencyService
	data            []*api.Product
}

func NewProductService(currencyService *CurrencyService) *ProductService {
	var productsList = data.LoadProductsFromJson("./data/products.json")
	return &ProductService{currencyService, productsList}
}

func (productService *ProductService) GetProducts() (api.Products, error) {
	rate, err := productService.CurrencyService.GetRate()
	if err != nil {
		return productService.data, &utils.GrpcServiceErr{Err: err.Error()}
	}

	ratedProductsList := api.Products{}
	for _, product := range productService.data {
		ratedProduct := *product
		ratedProduct.Price *= rate
		ratedProductsList = append(ratedProductsList, &ratedProduct)
	}

	return ratedProductsList, nil
}

func (productService *ProductService) GetProduct(id int) (*api.Product, error) {
	index := productService.findIndexByProductID(id)
	if index == -1 {
		return nil, &utils.ProductNotFoundErr{Err: fmt.Sprintf("id=%d", id)}
	}

	rate, err := productService.CurrencyService.GetRate()
	if err != nil {
		return productService.data[index], &utils.GrpcServiceErr{Err: err.Error()}
	}

	ratedProduct := *productService.data[index]
	ratedProduct.Price *= rate

	return &ratedProduct, nil
}

func (productService *ProductService) AddProduct(product *api.Product) error {
	product.ID = productService.GetNextProductId()

	rate, err := productService.CurrencyService.GetRate()
	if err != nil {
		err = &utils.GrpcServiceErr{Err: err.Error()}
	} else {
		product.Price *= rate
	}

	productService.data = append(productService.data, product)
	return err
}

func (productService *ProductService) UpdateProduct(product *api.Product) error {
	index := productService.findIndexByProductID(product.ID)
	if index == -1 {
		return &utils.ProductNotFoundErr{Err: fmt.Sprintf("id=%d", product.ID)}
	}

	rate, err := productService.CurrencyService.GetRate()
	if err != nil {
		err = &utils.GrpcServiceErr{Err: err.Error()}
	} else {
		product.Price *= rate
	}

	productService.data[index] = product
	return err
}

func (productService *ProductService) DeleteProduct(id int) error {
	index := productService.findIndexByProductID(id)
	if index == -1 {
		return &utils.ProductNotFoundErr{Err: fmt.Sprintf("id=%d", id)}
	}

	productService.data = productService.deleteProductByIndex(index)
	return nil
}

func (productService *ProductService) GetNextProductId() int {
	if len(productService.data) == 0 {
		return 1
	}

	return productService.data[len(productService.data)-1].ID + 1
}

/*
*
findIndexByProductID finds the index of a product.
Returns -1 when product not found.
*/
func (productService *ProductService) findIndexByProductID(id int) int {
	for index, product := range productService.data {
		if product.ID == id {
			return index
		}
	}

	return -1
}

func (productService *ProductService) deleteProductByIndex(index int) []*api.Product {
	return append(productService.data[:index], productService.data[index+1:]...)
}
