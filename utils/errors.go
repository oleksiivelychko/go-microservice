package utils

import "fmt"

type ProductNotFoundErr struct {
	Err string
}

type GrpcServiceErr struct {
	Err string
}

func (e *ProductNotFoundErr) Error() string {
	return fmt.Sprintf("product not found: %s", e.Err)
}

func (e *GrpcServiceErr) Error() string {
	return fmt.Sprintf("gRPC service error: %s", e.Err)
}
