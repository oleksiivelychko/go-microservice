package main

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	gService "github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
	"github.com/oleksiivelychko/go-microservice/backends"
	"github.com/oleksiivelychko/go-microservice/handlers"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-microservice/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const fileStorePrefix = "/files/"
const fileStoreBasePath = "./public" + fileStorePrefix
const swaggerPath = "/sdk/swagger.yaml"

func main() {
	var addr = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	var grpcAddr = fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("GRPC_PORT"))

	var origins = []string{
		"http://" + addr,
	}

	hcLogger := utils.NewLogger()
	validation := utils.NewValidation()

	// max file size is 5MB
	storage, err := backends.NewLocal(fileStoreBasePath, 1024*1000*5)
	if err != nil {
		hcLogger.Error("unable to create storage", "error", err)
		os.Exit(1)
	}

	grpcConnection, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		hcLogger.Error("unable to connect to gRPC server", "error", err)
	}
	defer grpcConnection.Close()

	currencyClient := gService.NewCurrencyClient(grpcConnection)
	currencyService := service.NewCurrencyService(currencyClient, "USD")
	productService := service.NewProductService(currencyService)

	productHandler := handlers.NewProductHandler(hcLogger, validation, productService)
	fileHandler := handlers.NewFileHandler(storage, hcLogger)
	multipartHandler := handlers.NewMultipartHandler(hcLogger, validation, storage, productService)
	gzipHandler := handlers.NewGzipHandler(hcLogger)

	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", productHandler.GetAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.GetOne)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", productHandler.CreateProduct)
	postRouter.Use(productHandler.MiddlewareProductValidation)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.UpdateProduct)
	putRouter.Use(productHandler.MiddlewareProductValidation)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", productHandler.DeleteProduct)

	// GET/POST file handling
	regex := fileStorePrefix + "{id:[0-9]+}/{filename:[a-zA-Z]+\\.(?:png|jpe?g)}"
	postFileRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postFileRouter.HandleFunc(regex, fileHandler.ServeHTTP)
	getRouter.Handle(regex, http.StripPrefix(fileStorePrefix, http.FileServer(http.Dir(fileStoreBasePath))))
	getRouter.Use(gzipHandler.MiddlewareGzip)

	// Multipart Form data processing
	postMultipartFormRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postMultipartFormRouter.HandleFunc("/products-form", multipartHandler.ProcessForm)

	opts := middleware.RedocOpts{SpecURL: swaggerPath}
	apiHandler := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", apiHandler)
	getRouter.Handle(swaggerPath, http.FileServer(http.Dir("./")))

	// Cross-Origin Resource Sharing
	goHandler := gohandlers.CORS(gohandlers.AllowedOrigins(origins))

	server := &http.Server{
		Addr:         addr,
		Handler:      goHandler(serveMux),
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

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_ = server.Shutdown(ctx)
}
