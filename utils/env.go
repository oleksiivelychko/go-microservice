package utils

import (
	"fmt"
	"os"
)

func GetServerAddr() (string, string) {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = ServerName
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = ServerPort
	}

	portGRPC, ok := os.LookupEnv("PORT_GRPC")
	if !ok {
		portGRPC = ServerPortGRPC
	}

	return fmt.Sprintf("%s:%s", host, port), fmt.Sprintf("%s:%s", host, portGRPC)
}
