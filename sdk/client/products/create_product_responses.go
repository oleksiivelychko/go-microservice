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

// CreateProductReader is a Reader for the CreateProduct structure.
type CreateProductReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateProductReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateProductCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewCreateProductBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewCreateProductUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewCreateProductCreated creates a CreateProductCreated with default headers values
func NewCreateProductCreated() *CreateProductCreated {
	return &CreateProductCreated{}
}

/*
CreateProductCreated describes a response with status code 201, with default header values.

Data structure representing a single product.
*/
type CreateProductCreated struct {
	Payload *models.Product
}

// IsSuccess returns true when this create product created response has a 2xx status code
func (o *CreateProductCreated) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this create product created response has a 3xx status code
func (o *CreateProductCreated) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create product created response has a 4xx status code
func (o *CreateProductCreated) IsClientError() bool {
	return false
}

// IsServerError returns true when this create product created response has a 5xx status code
func (o *CreateProductCreated) IsServerError() bool {
	return false
}

// IsCode returns true when this create product created response a status code equal to that given
func (o *CreateProductCreated) IsCode(code int) bool {
	return code == 201
}

func (o *CreateProductCreated) Error() string {
	return fmt.Sprintf("[POST /products][%d] createProductCreated  %+v", 201, o.Payload)
}

func (o *CreateProductCreated) String() string {
	return fmt.Sprintf("[POST /products][%d] createProductCreated  %+v", 201, o.Payload)
}

func (o *CreateProductCreated) GetPayload() *models.Product {
	return o.Payload
}

func (o *CreateProductCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Product)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateProductBadRequest creates a CreateProductBadRequest with default headers values
func NewCreateProductBadRequest() *CreateProductBadRequest {
	return &CreateProductBadRequest{}
}

/*
CreateProductBadRequest describes a response with status code 400, with default header values.

gRPC service request error message.
*/
type CreateProductBadRequest struct {
	Payload *models.GrpcError
}

// IsSuccess returns true when this create product bad request response has a 2xx status code
func (o *CreateProductBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create product bad request response has a 3xx status code
func (o *CreateProductBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create product bad request response has a 4xx status code
func (o *CreateProductBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this create product bad request response has a 5xx status code
func (o *CreateProductBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this create product bad request response a status code equal to that given
func (o *CreateProductBadRequest) IsCode(code int) bool {
	return code == 400
}

func (o *CreateProductBadRequest) Error() string {
	return fmt.Sprintf("[POST /products][%d] createProductBadRequest  %+v", 400, o.Payload)
}

func (o *CreateProductBadRequest) String() string {
	return fmt.Sprintf("[POST /products][%d] createProductBadRequest  %+v", 400, o.Payload)
}

func (o *CreateProductBadRequest) GetPayload() *models.GrpcError {
	return o.Payload
}

func (o *CreateProductBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateProductUnprocessableEntity creates a CreateProductUnprocessableEntity with default headers values
func NewCreateProductUnprocessableEntity() *CreateProductUnprocessableEntity {
	return &CreateProductUnprocessableEntity{}
}

/*
CreateProductUnprocessableEntity describes a response with status code 422, with default header values.

Validation errors defined as an array of strings.
*/
type CreateProductUnprocessableEntity struct {
	Payload *models.ValidationErrors
}

// IsSuccess returns true when this create product unprocessable entity response has a 2xx status code
func (o *CreateProductUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this create product unprocessable entity response has a 3xx status code
func (o *CreateProductUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this create product unprocessable entity response has a 4xx status code
func (o *CreateProductUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this create product unprocessable entity response has a 5xx status code
func (o *CreateProductUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this create product unprocessable entity response a status code equal to that given
func (o *CreateProductUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

func (o *CreateProductUnprocessableEntity) Error() string {
	return fmt.Sprintf("[POST /products][%d] createProductUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *CreateProductUnprocessableEntity) String() string {
	return fmt.Sprintf("[POST /products][%d] createProductUnprocessableEntity  %+v", 422, o.Payload)
}

func (o *CreateProductUnprocessableEntity) GetPayload() *models.ValidationErrors {
	return o.Payload
}

func (o *CreateProductUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ValidationErrors)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
