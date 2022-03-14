package utils

import "os"

func GetAddr() string {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "9090"
	}

	return host + ":" + port
}
