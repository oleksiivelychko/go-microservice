package main

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"github.com/oleksiivelychko/go-microservice/handler"
	"github.com/oleksiivelychko/go-microservice/handlers"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-microservice/utils"
	"github.com/oleksiivelychko/go-utils/logger"
	"github.com/oleksiivelychko/go-utils/server"
	"github.com/oleksiivelychko/go-utils/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	serverAddr, grpcServerAddr := utils.GetServerAddr()

	hcLogger := logger.NewHashicorp("go-microservice")
	validation := utils.NewValidation()

	localStorage, err := storage.NewLocal(utils.LocalStoragePath, utils.MaxFileSize5MB)
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
	currencyService := service.NewCurrencyService(hcLogger, exchangerClient, utils.DefaultCurrency)
	productService := service.NewProductService(currencyService, utils.LocalDataPath)

	productHandler := handler.NewProductHandler(hcLogger, validation, productService)
	fileHandler := handlers.NewFileHandler(localStorage, hcLogger)
	multipartHandler := handlers.NewMultipartHandler(hcLogger, validation, localStorage, productService)
	gzipHandler := handlers.NewHandlerGZIP(hcLogger)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.GetAll)
	getRouter.HandleFunc("/products", productHandler.GetAll).Queries(utils.CurrencyQueryParam, utils.CurrencyRegex)
	getRouter.HandleFunc(utils.ProductURL, productHandler.GetOne)
	getRouter.HandleFunc(utils.ProductURL, productHandler.GetOne).Queries(utils.CurrencyQueryParam, utils.CurrencyRegex)
	getRouter.Use(productHandler.MiddlewareCurrency)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", productHandler.CreateProduct)
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
	postMultipartFormRouter.HandleFunc("/products-form", multipartHandler.ProcessForm)

	swaggerUIOpts := middleware.SwaggerUIOpts{Path: utils.SwaggerURL, SpecURL: utils.SwaggerYAML}
	swaggerUI := middleware.SwaggerUI(swaggerUIOpts, nil)
	getRouter.Handle(utils.SwaggerURL, swaggerUI)

	redocOpts := middleware.RedocOpts{Path: utils.RedocURL, SpecURL: utils.SwaggerYAML}
	redoc := middleware.Redoc(redocOpts, nil)
	getRouter.Handle(utils.RedocURL, redoc)

	getRouter.Handle(utils.SwaggerYAML, http.FileServer(http.Dir("./")))

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
