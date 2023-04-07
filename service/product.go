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

func NewProductService(currencyService *CurrencyService, filePath string) *ProductService {
	var products = data.LoadProductsFromJson(filePath)
	return &ProductService{currencyService, products}
}

func (service *ProductService) GetProducts() (api.Products, *errors.GRPCServiceError) {
	rate, grpcErr := service.CurrencyService.GetRate()
	if grpcErr != nil {
		return service.products, grpcErr
	}

	ratedProductsList := api.Products{}
	for _, product := range service.products {
		ratedProduct := *product
		ratedProduct.Price *= rate
		ratedProductsList = append(ratedProductsList, &ratedProduct)
	}

	return ratedProductsList, nil
}

func (service *ProductService) GetProduct(id int) (*api.Product, error) {
	index := service.findIndexByProductID(id)
	if index == -1 {
		return nil, &errors.ProductNotFoundError{ID: id}
	}

	rate, grpcErr := service.CurrencyService.GetRate()
	if grpcErr != nil {
		return service.products[index], grpcErr
	}

	ratedProduct := *service.products[index]
	ratedProduct.Price *= rate

	return &ratedProduct, nil
}

func (service *ProductService) AddProduct(product *api.Product) *errors.GRPCServiceError {
	product.ID = service.GetNextProductID()

	rate, grpcErr := service.CurrencyService.GetRate()
	if grpcErr != nil {
		return grpcErr
	}

	product.Price *= rate
	service.products = append(service.products, product)

	return nil
}

func (service *ProductService) UpdateProduct(product *api.Product) error {
	index := service.findIndexByProductID(product.ID)
	if index == -1 {
		return &errors.ProductNotFoundError{ID: product.ID}
	}

	rate, grpcErr := service.CurrencyService.GetRate()
	if grpcErr != nil {
		return grpcErr
	}

	product.Price *= rate
	service.products[index] = product

	return nil
}

func (service *ProductService) DeleteProduct(id int) *errors.ProductNotFoundError {
	index := service.findIndexByProductID(id)
	if index == -1 {
		return &errors.ProductNotFoundError{ID: id}
	}

	service.products = service.deleteProductByIndex(index)
	return nil
}

func (service *ProductService) GetNextProductID() int {
	if len(service.products) == 0 {
		return 1
	}

	return service.products[len(service.products)-1].ID + 1
}

/*
*
findIndexByProductID finds the index of a product. Returns -1 when product not found.
*/
func (service *ProductService) findIndexByProductID(id int) int {
	for index, product := range service.products {
		if product.ID == id {
			return index
		}
	}

	return -1
}

func (service *ProductService) deleteProductByIndex(index int) []*api.Product {
	return append(service.products[:index], service.products[index+1:]...)
}
