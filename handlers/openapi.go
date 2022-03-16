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

// NOTE: Types defined here are purely for documentation purposes these types are not used by any of the handlers

// Generic error message returned as a string.
// swagger:response errorResponse
type errResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Validation errors defined as an array of strings.
// swagger:response errorValidation
type errValidationWrapper struct {
	// Collection of the errors
	// in: body
	Body ValidationError
}

// Data structure representing a list of product
// swagger:response productsResponse
type productsResponseWrapper struct {
	// A list of all products
	// in: body
	Body []api.Product
}

// Data structure representing a single product
// swagger:response productResponse
type productResponseWrapper struct {
	// Newly created product
	// in: body
	Body api.Product
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct{}

// swagger:parameters updateProduct createProduct
type productParamsWrapper struct {
	// Product data structure to Update or Create.
	// Note: the id field is ignored by update and create operations
	// in: body
	// required: true
	Body api.Product
}

// swagger:parameters getProduct deleteProduct
type productIDParamWrapper struct {
	// The id of the product for which the operation relates
	// in: path
	// required: true
	ID int `json:"id"`
}
