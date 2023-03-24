package errors

import "fmt"

// GenericError is a generic error message returned by a server.
type GenericError struct {
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation error messages.
type ValidationErrors struct {
	Messages []string `json:"messages"`
}

type ProductNotFoundError struct {
	ID int `json:"id"`
}

type GRPCServiceError struct {
	Message string `json:"message"`
}

func (err *ProductNotFoundError) Error() string {
	return fmt.Sprintf("product %d not found", err.ID)
}

func (err *GRPCServiceError) Error() string {
	return fmt.Sprintf("gRPC service error: %s", err.Message)
}
