package utils

import "fmt"

type ProductNotFoundErr struct {
	Err string
}

type GrpcServiceRequestErr struct {
	Err string
}

func (e *ProductNotFoundErr) Error() string {
	return fmt.Sprintf("product not found: %s", e.Err)
}

func (e *GrpcServiceRequestErr) Error() string {
	return fmt.Sprintf("unable to make request to gRPC service. %s", e.Err)
}
