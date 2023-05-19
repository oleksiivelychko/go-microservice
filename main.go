package main

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oleksiivelychko/go-grpc-service/logger"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"github.com/oleksiivelychko/go-microservice/handlers"
	"github.com/oleksiivelychko/go-microservice/handlers/product"
	"github.com/oleksiivelychko/go-microservice/server"
	"github.com/oleksiivelychko/go-microservice/services"
	"github.com/oleksiivelychko/go-microservice/storage"
	"github.com/oleksiivelychko/go-microservice/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const defaultCurrency = "USD"
const localDataPath = "./public/data/products.json"
const localStorageBasePath = "/files/"
const localStoragePath = "./public" + localStorageBasePath
const maxFileSize5MB = 1024 * 1000 * 5
const RedocURL = "/redoc"
const SwaggerURL = "/swagger"
const SwaggerYAML = "/sdk/swagger.yaml"

func main() {
	serverAddr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	grpcServerAddr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT_GRPC"))

	log := logger.New()

	validate, err := validation.New()
	if err != nil {
		log.Error("unable to create validator: %s", err)
		os.Exit(1)
	}

	localStorage, err := storage.New(localStoragePath, maxFileSize5MB)
	if err != nil {
		log.Error("unable to create local storage: %s", err)
		os.Exit(1)
	}

	grpcConnection, err := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("unable to connect to gRPC server: %s", err)
	}
	defer grpcConnection.Close()

	exchangerClient := grpcservice.NewExchangerClient(grpcConnection)
	currencyService := services.NewCurrency(exchangerClient, defaultCurrency, log)
	productService := services.NewProduct(currencyService, localDataPath)

	productHandler := product.NewHandler(validate, productService, log)
	fileHandler := handlers.NewFile(localStorage, log)
	multipartHandler := handlers.NewMultipart(validate, localStorage, productService, log)
	gzipHandler := handlers.NewGZIP(log)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.GetAll)
	getRouter.HandleFunc("/products", productHandler.GetAll).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.GetOne)
	getRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.GetOne).Queries("currency", "{[A-Z]{3}}")
	getRouter.Use(productHandler.MiddlewareCurrency)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", productHandler.CreateProduct)
	postRouter.Use(productHandler.MiddlewareValidation)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.UpdateProduct)
	putRouter.Use(productHandler.MiddlewareValidation)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.DeleteProduct)

	// GET/POST file handling
	fileNameRegex := localStorageBasePath + "{id:[0-9]+}/{filename:[a-zA-Z]+\\.(?:png|jpe?g)}"
	postFileRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postFileRouter.HandleFunc(fileNameRegex, fileHandler.ServeHTTP)
	getRouter.Handle(fileNameRegex, http.StripPrefix(localStorageBasePath, http.FileServer(http.Dir(localStoragePath))))
	getRouter.Use(gzipHandler.Middleware)

	// multipart/form-data processing
	postMultipartFormRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postMultipartFormRouter.HandleFunc("/products-form", multipartHandler.ProcessForm)

	swaggerUIOpts := middleware.SwaggerUIOpts{Path: SwaggerURL, SpecURL: SwaggerYAML}
	swaggerUI := middleware.SwaggerUI(swaggerUIOpts, nil)
	getRouter.Handle(SwaggerURL, swaggerUI)

	redocOpts := middleware.RedocOpts{Path: RedocURL, SpecURL: SwaggerYAML}
	redoc := middleware.Redoc(redocOpts, nil)
	getRouter.Handle(RedocURL, redoc)

	getRouter.Handle(SwaggerYAML, http.FileServer(http.Dir("./")))

	// Cross-Origin Resource Sharing
	goHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{
		"http://" + serverAddr,
	}))

	httpServer := server.NewHTTP(serverAddr, goHandler(serveMux), log.GetErrorLogger())

	go func() {
		log.Info("starting server on %s", serverAddr)
		err = httpServer.ListenAndServe()
		if err != nil {
			log.Error("unable to start server: %s", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the httpServer
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	// block until a signal is received
	signalCh := <-signalChannel
	log.Info("received terminate, graceful shutdown; signal=%s", signalCh)

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// gracefully shutdown the httpServer, waiting max 30 seconds for current operations to complete
	httpServer.Shutdown(contextWithTimeout)
}
