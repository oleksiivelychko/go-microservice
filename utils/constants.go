package utils

const CurrencyQueryParam = "currency"
const CurrencyRegex = "{[A-Z]{3}}"
const DefaultCurrency = "USD"
const FormDataMaxMemory32MB = 128 * 1024
const LocalDataPath = "./data/products.json"
const LocalStorageBasePath = "/files/"
const LocalStoragePath = "./public" + LocalStorageBasePath
const MaxFileSize5MB = 1024 * 1000 * 5
const ProductFileURL = "{id:[0-9]+}/{filename:[a-zA-Z]+\\.(?:png|jpe?g)}"
const ProductURL = "/products/{id:[0-9]+}"
const RedocURL = "/redoc"
const SwaggerURL = "/swagger"
const SwaggerYAML = "/sdk/swagger.yaml"
