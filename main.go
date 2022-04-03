package main

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oleksiivelychko/go-helper/env"
	"github.com/oleksiivelychko/go-microservice/handlers"
	"github.com/oleksiivelychko/go-microservice/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	addr := env.GetAddr()
	var origins = []string{
		"http://" + addr,
	}

	l := log.New(os.Stdout, "go-microservice", log.LstdFlags)
	v := utils.NewValidation()

	h := handlers.NewProductHandler(l, v)
	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", h.GetAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", h.GetOne)

	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", h.CreateProduct)
	postRouter.Use(h.MiddlewareProductValidation)

	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products", h.UpdateProduct)
	putRouter.Use(h.MiddlewareProductValidation)

	deleteRouter := serveMux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", h.DeleteProduct)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	apiHandler := middleware.Redoc(opts, nil)
	getRouter.Handle("/docs", apiHandler)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// Cross-Origin Resource Sharing
	goHandler := gohandlers.CORS(gohandlers.AllowedOrigins(origins))

	server := &http.Server{
		Addr:         addr,
		Handler:      goHandler(serveMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Printf("Starting server on %s\n", addr)

		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	log.Println("Received terminate, graceful shutdown", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_ = server.Shutdown(ctx)
}
