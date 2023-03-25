/*
Documentation of microservice API

	Schemes:
	  http
	Host: localhost
	BasePath: /
	Version: 1.0.0
	Consumes:
	- application/json
	Produces:
	- application/json

swagger:meta
*/
package errors

import (
	"github.com/oleksiivelychko/go-microservice/api"
)

// NOTE: Types defined here are purely for documentation purposes, these types are not used by any of the handlers.

// swagger:response errorResponse
type errorResponse struct {
	// in: body
	Body GenericError
}

// Validation errors are an array of strings.
// swagger:response validationErrorsResponse
type validationErrorsResponse struct {
	// in: body
	Body ValidationErrors
}

// swagger:response grpcErrorResponse
type grpcErrorResponse struct {
	// in: body
	Body GRPCServiceError
}

// swagger:response notFoundResponse
type notFoundResponse struct {
	// in: body
	Body ProductNotFoundError
}

// Data structure is representing a list of products.
// swagger:response productsResponse
type productsResponse struct {
	// in: body
	Body []api.Product
}

// Data structure is representing a single product.
// swagger:response productResponse
type productResponse struct {
	// in: body
	Body api.Product
}

// Empty response has no data.
// swagger:response noContentResponse
type noContentResponse struct{}

// Send product data as part of HTTP request (ID field would be ignored).
// swagger:parameters createProduct
type productPostRequest struct {
	// in: body
	// required: true
	Body api.Product
}

// Request product by ID parameter from URL and send data in body.
// swagger:parameters updateProduct
type productPutRequest struct {
	// in: path
	// required: true
	ID int `json:"id"`
	// in: body
	// required: true
	Body api.Product
}

// Request product by ID parameter.
// swagger:parameters getProduct deleteProduct
type productGetRequest struct {
	// in: path
	// required: true
	ID int `json:"id"`
}

// Hand over currency optional parameter.
// swagger:parameters getProducts getProduct
type productQueryParameters struct {
	// in: query
	// required: false
	Currency string
}
