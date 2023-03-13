package main

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	grpcService "github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	"github.com/oleksiivelychko/go-microservice/handlers"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-microservice/utils"
	localStorage "github.com/oleksiivelychko/go-utils/local_storage"
	"github.com/oleksiivelychko/go-utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const fileStorePrefix = "/files/"
const fileStoreBasePath = "./public" + fileStorePrefix

func main() {
	var addr = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	var grpcAddr = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("GRPC_PORT"))

	hcLogger := logger.NewLogger("go-microservice")
	validation := utils.NewValidation()

	storage, err := localStorage.NewLocalStorage(fileStoreBasePath, 1024*1000*5) // max file size is 5Mb
	if err != nil {
		hcLogger.Error("unable to create localStorage", "error", err)
		os.Exit(1)
	}

	grpcConnection, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hcLogger.Error("unable to connect to gRPC server", "error", err)
	}
	defer grpcConnection.Close()

	currencyClient := grpcService.NewCurrencyClient(grpcConnection)
	currencyService := service.NewCurrencyService(hcLogger, currencyClient, "USD")
	productService := service.NewProductService(currencyService)

	productHandler := handlers.NewProductHandler(hcLogger, validation, productService)
	fileHandler := handlers.NewFileHandler(storage, hcLogger)
	multipartHandler := handlers.NewMultipartHandler(hcLogger, validation, storage, productService)
	gzipHandler := handlers.NewGzipHandler(hcLogger)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.GetAll)
	getRouter.HandleFunc("/products", productHandler.GetAll).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.GetOne)
	getRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.GetOne).Queries("currency", "{[A-Z]{3}}")
	getRouter.Use(productHandler.MiddlewareProductCurrency)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", productHandler.CreateProduct)
	postRouter.Use(productHandler.MiddlewareProductValidation)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.UpdateProduct)
	putRouter.Use(productHandler.MiddlewareProductValidation)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.DeleteProduct)

	// GET/POST file handling
	var fileNameRegex = fileStorePrefix + "{id:[0-9]+}/{filename:[a-zA-Z]+\\.(?:png|jpe?g)}"
	postFileRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postFileRouter.HandleFunc(fileNameRegex, fileHandler.ServeHTTP)
	getRouter.Handle(fileNameRegex, http.StripPrefix(fileStorePrefix, http.FileServer(http.Dir(fileStoreBasePath))))
	getRouter.Use(gzipHandler.MiddlewareGzip)

	// multipart/form-data processing
	postMultipartFormRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postMultipartFormRouter.HandleFunc("/products-form", multipartHandler.ProcessForm)

	var swaggerPath = "/sdk/swagger.yaml"
	redocOpts := middleware.RedocOpts{SpecURL: swaggerPath}
	apiHandler := middleware.Redoc(redocOpts, nil)
	getRouter.Handle("/docs", apiHandler)
	getRouter.Handle(swaggerPath, http.FileServer(http.Dir("./")))

	// Cross-Origin Resource Sharing
	handler := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{
		"http://" + addr,
	}))

	server := &http.Server{
		Addr:         addr,
		Handler:      handler(serveMux),
		ErrorLog:     hcLogger.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true}),
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
		ReadTimeout:  10 * time.Second,  // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
	}

	go func() {
		hcLogger.Info("starting server", "listening", addr)
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
	sig := <-signalChannel
	hcLogger.Info("received terminate, graceful shutdown", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	server.Shutdown(ctx)
}
