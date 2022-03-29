// Code generated by go-swagger; DO NOT EDIT.

package products

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/oleksiivelychko/go-microservice/sdk/models"
)

// UpdateProductReader is a Reader for the UpdateProduct structure.
type UpdateProductReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateProductReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewUpdateProductCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewUpdateProductNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewUpdateProductUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 501:
		result := NewUpdateProductNotImplemented()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewUpdateProductCreated creates a UpdateProductCreated with default headers values
func NewUpdateProductCreated() *UpdateProductCreated {
	return &UpdateProductCreated{}
}

/* UpdateProductCreated describes a response with status code 201, with default header values.

Data structure representing a single product
*/
type UpdateProductCreated struct {
	Payload *models.Product
}

func (o *UpdateProductCreated) Error() string {
	return fmt.Sprintf("[PUT /products/{id}][%d] updateProductCreated  %+v", 201, o.Payload)
}
func (o *UpdateProductCreated) GetPayload() *models.Product {
	return o.Payload
}

func (o *UpdateProductCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Product)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateProductNotFound creates a UpdateProductNotFound with default headers values
func NewUpdateProductNotFound() *UpdateProductNotFound {
	return &UpdateProductNotFound{}
}

/* UpdateProductNotFound describes a response with status code 404, with default header values.

Generic error message returned as a string.
*/
type UpdateProductNotFound struct {
	Payload *models.GenericError
}

func (o *UpdateProductNotFound) Error() string {
	return fmt.Sprintf("[PUT /products/{id}][%d] updateProductNotFound  %+v", 404, o.Payload)
}
func (o *UpdateProductNotFound) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *UpdateProductNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateProductUnprocessableEntity creates a UpdateProductUnprocessableEntity with default headers values
func NewUpdateProductUnprocessableEntity() *UpdateProductUnprocessableEntity {
	return &UpdateProductUnprocessableEntity{}
}

/* UpdateProductUnprocessableEntity describes a response with status code 422, with default header values.

Validation errors defined as an array of strings.
*/
type UpdateProductUnprocessableEntity struct {
	Payload *models.ValidationError
}

func (o *UpdateProductUnprocessableEntity) Error() string {
	return fmt.Sprintf("[PUT /products/{id}][%d] updateProductUnprocessableEntity  %+v", 422, o.Payload)
}
func (o *UpdateProductUnprocessableEntity) GetPayload() *models.ValidationError {
	return o.Payload
}

func (o *UpdateProductUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ValidationError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateProductNotImplemented creates a UpdateProductNotImplemented with default headers values
func NewUpdateProductNotImplemented() *UpdateProductNotImplemented {
	return &UpdateProductNotImplemented{}
}

/* UpdateProductNotImplemented describes a response with status code 501, with default header values.

Generic error message returned as a string.
*/
type UpdateProductNotImplemented struct {
	Payload *models.GenericError
}

func (o *UpdateProductNotImplemented) Error() string {
	return fmt.Sprintf("[PUT /products/{id}][%d] updateProductNotImplemented  %+v", 501, o.Payload)
}
func (o *UpdateProductNotImplemented) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *UpdateProductNotImplemented) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
