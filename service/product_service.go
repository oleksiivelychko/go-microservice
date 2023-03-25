package service

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/data"
	"github.com/oleksiivelychko/go-microservice/errors"
)

type ProductService struct {
	CurrencyService *CurrencyService
	products        []*api.Product
}

func NewProductService(currencyService *CurrencyService, localDataPath string) *ProductService {
	var productsList = data.LoadProductsFromJson(localDataPath)
	return &ProductService{currencyService, productsList}
}

func (productService *ProductService) GetProducts() (api.Products, *errors.GRPCServiceError) {
	rate, grpcErr := productService.CurrencyService.GetRate()
	if grpcErr != nil {
		return productService.products, grpcErr
	}

	ratedProductsList := api.Products{}
	for _, product := range productService.products {
		ratedProduct := *product
		ratedProduct.Price *= rate
		ratedProductsList = append(ratedProductsList, &ratedProduct)
	}

	return ratedProductsList, nil
}

func (productService *ProductService) GetProduct(id int) (*api.Product, error) {
	index := productService.findIndexByProductID(id)
	if index == -1 {
		return nil, &errors.ProductNotFoundError{ID: id}
	}

	rate, grpcErr := productService.CurrencyService.GetRate()
	if grpcErr != nil {
		return productService.products[index], grpcErr
	}

	ratedProduct := *productService.products[index]
	ratedProduct.Price *= rate

	return &ratedProduct, nil
}

func (productService *ProductService) AddProduct(product *api.Product) *errors.GRPCServiceError {
	product.ID = productService.GetNextProductID()

	rate, grpcErr := productService.CurrencyService.GetRate()
	if grpcErr != nil {
		return grpcErr
	}

	product.Price *= rate
	productService.products = append(productService.products, product)

	return nil
}

func (productService *ProductService) UpdateProduct(product *api.Product) error {
	index := productService.findIndexByProductID(product.ID)
	if index == -1 {
		return &errors.ProductNotFoundError{ID: product.ID}
	}

	rate, grpcErr := productService.CurrencyService.GetRate()
	if grpcErr != nil {
		return grpcErr
	}

	product.Price *= rate
	productService.products[index] = product

	return nil
}

func (productService *ProductService) DeleteProduct(id int) *errors.ProductNotFoundError {
	index := productService.findIndexByProductID(id)
	if index == -1 {
		return &errors.ProductNotFoundError{ID: id}
	}

	productService.products = productService.deleteProductByIndex(index)
	return nil
}

func (productService *ProductService) GetNextProductID() int {
	if len(productService.products) == 0 {
		return 1
	}

	return productService.products[len(productService.products)-1].ID + 1
}

/*
*
findIndexByProductID finds the index of a product. Returns -1 when product not found.
*/
func (productService *ProductService) findIndexByProductID(id int) int {
	for index, product := range productService.products {
		if product.ID == id {
			return index
		}
	}

	return -1
}

func (productService *ProductService) deleteProductByIndex(index int) []*api.Product {
	return append(productService.products[:index], productService.products[index+1:]...)
}
