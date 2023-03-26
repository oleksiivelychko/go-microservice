package utils

const CurrencyQueryParam = "currency"
const CurrencyRegex = "{[A-Z]{3}}"
const DefaultCurrency = "USD"
const LocalDataPath = "./data/products.json"
const LocalStorageBasePath = "/files/"
const LocalStoragePath = "./public" + LocalStorageBasePath
const MaxFileSize5MB = 1024 * 1000 * 5
const ProductFileURL = "{id:[0-9]+}/{filename:[a-zA-Z]+\\.(?:png|jpe?g)}"
const ProductURL = "/products/{id:[0-9]+}"
const ProductsFormURL = "/products-form"
const ProductsURL = "/products"
const RedocURL = "/redoc"
const ServerName = "localhost"
const ServerPort = "9090"
const ServerPortGRPC = "9091"
const ServerAddr = ServerName + "/" + ServerPort
const SwaggerURL = "/swagger"
const SwaggerYAML = "/sdk/swagger.yaml"
