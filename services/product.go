package services

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/data"
	"github.com/oleksiivelychko/go-microservice/errors"
)

type Product struct {
	CurrencyService *Currency
	products        []*api.Product
}

func NewProduct(currencyService *Currency, filePath string) *Product {
	var products = data.LoadProductsFromJSON(filePath)
	return &Product{currencyService, products}
}

func (service *Product) GetProducts() (api.Products, *errors.GRPCServiceError) {
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

func (service *Product) GetProduct(id int) (*api.Product, error) {
	idx := service.findIndexByProductID(id)
	if idx == -1 {
		return nil, &errors.ProductNotFoundError{ID: id}
	}

	rate, grpcErr := service.CurrencyService.GetRate()
	if grpcErr != nil {
		return service.products[idx], grpcErr
	}

	ratedProduct := *service.products[idx]
	ratedProduct.Price *= rate

	return &ratedProduct, nil
}

func (service *Product) AddProduct(product *api.Product) *errors.GRPCServiceError {
	product.ID = service.GetNextProductID()

	rate, grpcErr := service.CurrencyService.GetRate()
	if grpcErr != nil {
		return grpcErr
	}

	product.Price *= rate
	service.products = append(service.products, product)

	return nil
}

func (service *Product) UpdateProduct(product *api.Product) error {
	idx := service.findIndexByProductID(product.ID)
	if idx == -1 {
		return &errors.ProductNotFoundError{ID: product.ID}
	}

	rate, grpcErr := service.CurrencyService.GetRate()
	if grpcErr != nil {
		return grpcErr
	}

	product.Price *= rate
	service.products[idx] = product

	return nil
}

func (service *Product) DeleteProduct(id int) *errors.ProductNotFoundError {
	idx := service.findIndexByProductID(id)
	if idx == -1 {
		return &errors.ProductNotFoundError{ID: id}
	}

	service.products = service.deleteProductByIndex(idx)
	return nil
}

func (service *Product) GetNextProductID() int {
	if len(service.products) == 0 {
		return 1
	}

	return service.products[len(service.products)-1].ID + 1
}

/*
*
findIndexByProductID finds the index of a product. Returns -1 when product not found.
*/
func (service *Product) findIndexByProductID(id int) int {
	for idx, product := range service.products {
		if product.ID == id {
			return idx
		}
	}

	return -1
}

func (service *Product) deleteProductByIndex(idx int) []*api.Product {
	return append(service.products[:idx], service.products[idx+1:]...)
}
