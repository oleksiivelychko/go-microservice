package main

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	"github.com/oleksiivelychko/go-microservice/handlers"
	"github.com/oleksiivelychko/go-microservice/product_handler"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-microservice/utils"
	"github.com/oleksiivelychko/go-utils/local_storage"
	"github.com/oleksiivelychko/go-utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var hostAddr = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	var grpcAddr = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("GRPC_PORT"))

	hcLogger := logger.NewLogger("go-microservice")
	validation := utils.NewValidation()

	localStorage, err := local_storage.NewLocalStorage(utils.LocalStoragePath, utils.MaxFileSize5MB)
	if err != nil {
		hcLogger.Error("unable to create local storage", "error", err)
		os.Exit(1)
	}

	grpcConnection, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hcLogger.Error("unable to connect to gRPC server", "error", err)
	}
	defer grpcConnection.Close()

	currencyClient := grpc_service.NewCurrencyClient(grpcConnection)
	currencyService := service.NewCurrencyService(hcLogger, currencyClient, utils.DefaultCurrency)
	productService := service.NewProductService(currencyService, utils.LocalDataPath)

	productHandler := product_handler.NewProductHandler(hcLogger, validation, productService)
	fileHandler := handlers.NewFileHandler(localStorage, hcLogger)
	multipartHandler := handlers.NewMultipartHandler(hcLogger, validation, localStorage, productService)
	gzipHandler := handlers.NewGzipHandler(hcLogger)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc(utils.ProductsURL, productHandler.GetAll)
	getRouter.HandleFunc(utils.ProductsURL, productHandler.GetAll).Queries(utils.CurrencyQueryParam, utils.CurrencyRegex)
	getRouter.HandleFunc(utils.ProductURL, productHandler.GetOne)
	getRouter.HandleFunc(utils.ProductURL, productHandler.GetOne).Queries(utils.CurrencyQueryParam, utils.CurrencyRegex)
	getRouter.Use(productHandler.MiddlewareCurrency)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc(utils.ProductsURL, productHandler.CreateProduct)
	postRouter.Use(productHandler.MiddlewareValidation)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc(utils.ProductURL, productHandler.UpdateProduct)
	putRouter.Use(productHandler.MiddlewareValidation)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc(utils.ProductURL, productHandler.DeleteProduct)

	// GET/POST file handling
	var fileNameRegex = utils.LocalStorageBasePath + utils.ProductFileURL
	postFileRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postFileRouter.HandleFunc(fileNameRegex, fileHandler.ServeHTTP)
	getRouter.Handle(fileNameRegex, http.StripPrefix(
		utils.LocalStorageBasePath, http.FileServer(http.Dir(utils.LocalStoragePath)),
	))
	getRouter.Use(gzipHandler.Middleware)

	// multipart/form-data processing
	postMultipartFormRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postMultipartFormRouter.HandleFunc(utils.ProductsFormURL, multipartHandler.ProcessForm)

	swaggerUIOpts := middleware.SwaggerUIOpts{SpecURL: utils.SwaggerYAML}
	swaggerUI := middleware.SwaggerUI(swaggerUIOpts, nil)
	getRouter.Handle(utils.SwaggerURL, swaggerUI)

	redocOpts := middleware.RedocOpts{SpecURL: utils.SwaggerYAML}
	redoc := middleware.Redoc(redocOpts, nil)
	getRouter.Handle(utils.RedocURL, redoc)

	getRouter.Handle(utils.SwaggerYAML, http.FileServer(http.Dir("./")))

	// Cross-Origin Resource Sharing
	handler := gorillahandlers.CORS(gorillahandlers.AllowedOrigins([]string{
		"http://" + hostAddr,
	}))

	server := &http.Server{
		Addr:         hostAddr,
		Handler:      handler(serveMux),
		ErrorLog:     hcLogger.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true}),
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
		ReadTimeout:  10 * time.Second,  // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
	}

	go func() {
		hcLogger.Info("starting server", "listening", hostAddr)
		err = server.ListenAndServe()
		if err != nil {
			hcLogger.Error("unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	// block until a signal is received
	signalCh := <-signalChannel
	hcLogger.Info("received terminate, graceful shutdown", "signal", signalCh)

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	server.Shutdown(contextWithTimeout)
}
