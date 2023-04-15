package env

import (
	"fmt"
	"os"
)

const DefaultCurrency = "USD"
const FormDataMaxMemory32MB = 128 * 1024
const LocalDataPath = "./public/data/products.json"
const LocalStorageBasePath = "/files/"
const LocalStoragePath = "./public" + LocalStorageBasePath
const MaxFileSize5MB = 1024 * 1000 * 5
const ProductFileURL = "{id:[0-9]+}/{filename:[a-zA-Z]+\\.(?:png|jpe?g)}"
const RedocURL = "/redoc"
const SwaggerURL = "/swagger"
const SwaggerYAML = "/sdk/swagger.yaml"

func ServerAddr() (string, string) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	portGRPC := os.Getenv("PORT_GRPC")

	return fmt.Sprintf("%s:%s", host, port), fmt.Sprintf("%s:%s", host, portGRPC)
}
