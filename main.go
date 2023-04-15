package main

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"github.com/oleksiivelychko/go-microservice/env"
	"github.com/oleksiivelychko/go-microservice/handler"
	"github.com/oleksiivelychko/go-microservice/handlers"
	"github.com/oleksiivelychko/go-microservice/logger"
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

func main() {
	serverAddr, grpcServerAddr := env.ServerAddr()
	hcLogger := logger.New("go-microservice", "DEBUG")

	validate, err := validation.New()
	if err != nil {
		hcLogger.Error("unable to create validator", "error", err)
		os.Exit(1)
	}

	localStorage, err := storage.New(env.LocalStoragePath, env.MaxFileSize5MB)
	if err != nil {
		hcLogger.Error("unable to create local storage", "error", err)
		os.Exit(1)
	}

	grpcConnection, err := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hcLogger.Error("unable to connect to gRPC server", "error", err)
	}
	defer grpcConnection.Close()

	exchangerClient := grpcservice.NewExchangerClient(grpcConnection)
	currencyService := services.NewCurrency(hcLogger, exchangerClient, env.DefaultCurrency)
	productService := services.NewProduct(currencyService, env.LocalDataPath)

	productHandler := handler.New(hcLogger, validate, productService)
	fileHandler := handlers.NewFile(localStorage, hcLogger)
	multipartHandler := handlers.NewMultipart(hcLogger, validate, localStorage, productService)
	gzipHandler := handlers.NewGZIP(hcLogger)

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
	var fileNameRegex = env.LocalStorageBasePath + env.ProductFileURL
	postFileRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postFileRouter.HandleFunc(fileNameRegex, fileHandler.ServeHTTP)
	getRouter.Handle(fileNameRegex, http.StripPrefix(
		env.LocalStorageBasePath, http.FileServer(http.Dir(env.LocalStoragePath)),
	))
	getRouter.Use(gzipHandler.Middleware)

	// multipart/form-data processing
	postMultipartFormRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postMultipartFormRouter.HandleFunc("/products-form", multipartHandler.ProcessForm)

	swaggerUIOpts := middleware.SwaggerUIOpts{Path: env.SwaggerURL, SpecURL: env.SwaggerYAML}
	swaggerUI := middleware.SwaggerUI(swaggerUIOpts, nil)
	getRouter.Handle(env.SwaggerURL, swaggerUI)

	redocOpts := middleware.RedocOpts{Path: env.RedocURL, SpecURL: env.SwaggerYAML}
	redoc := middleware.Redoc(redocOpts, nil)
	getRouter.Handle(env.RedocURL, redoc)

	getRouter.Handle(env.SwaggerYAML, http.FileServer(http.Dir("./")))

	// Cross-Origin Resource Sharing
	goHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{
		"http://" + serverAddr,
	}))

	httpServer := server.NewHTTP(serverAddr, goHandler(serveMux), hcLogger)

	go func() {
		hcLogger.Info("starting server", "listening", serverAddr)
		err = httpServer.ListenAndServe()
		if err != nil {
			hcLogger.Error("unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the httpServer
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	// block until a signal is received
	signalCh := <-signalChannel
	hcLogger.Info("received terminate, graceful shutdown", "signal", signalCh)

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// gracefully shutdown the httpServer, waiting max 30 seconds for current operations to complete
	httpServer.Shutdown(contextWithTimeout)
}
