package utils

import (
	"fmt"
	"os"
)

func GetServerAddr() (string, string) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	portGRPC := os.Getenv("PORT_GRPC")

	return fmt.Sprintf("%s:%s", host, port), fmt.Sprintf("%s:%s", host, portGRPC)
}
