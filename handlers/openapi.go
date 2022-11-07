/*
Documentation of microservice API

	Schemes:
	  http
	  https
	Host: localhost
	BasePath: /
	Version: 1.0.0
	Consumes:
	- application/json
	Produces:
	- application/json

swagger:meta
*/
package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
)

// NOTE: Types defined here are purely for documentation purposes, these types are not used by any of the handlers.

// Generic error message returned as a string.
// swagger:response errorResponse
type errorResponseWrapper struct {
	// in: body
	Body GenericError
}

// Validation errors defined as an array of strings.
// swagger:response validationErrorsResponse
type validationErrorsResponseWrapper struct {
	// in: body
	Body ValidationErrors
}

// swagger:response notFoundResponse
type notFoundResponseWrapper struct {
	// in: body
	Body NotFound
}

// gRPC service request error message.
// swagger:response grpcResponseWrapper
type grpcResponseWrapper struct {
	// in: body
	Body GrpcError
}

// Data structure representing a list of product.
// swagger:response productsResponse
type productsResponseWrapper struct {
	// in: body
	Body []api.Product
}

// Data structure representing a single product.
// swagger:response productResponse
type productResponseWrapper struct {
	// in: body
	Body api.Product
}

// Empty response has no data.
// swagger:response noContentResponse
type noContentResponseWrapper struct{}

// Send product data as part of HTTP request (ID field would be ignored).
// swagger:parameters createProduct
type productRequestBodyWrapper struct {
	// in: body
	// required: true
	Body api.Product
}

// Request product by ID parameter in URL and send data in body.
// swagger:parameters updateProduct
type productRequestIdBodyWrapper struct {
	// in: path
	// required: true
	ID int `json:"id"`
	// in: body
	// required: true
	Body api.Product
}

// Request product by ID parameter.
// swagger:parameters getProduct deleteProduct
type productRequestIdWrapper struct {
	// in: path
	// required: true
	ID int `json:"id"`
}

// Hand over currency optional parameter.
// swagger:parameters getProducts getProduct
type productQueryCurrencyWrapper struct {
	// in: query
	// required: false
	Currency string
}
