package handlers

import (
	"github.com/gorilla/mux"
	"github.com/oleksiivelychko/go-microservice/api"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products", r.URL.Path)

	lp := api.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products", r.URL.Path)

	product := &api.Product{}

	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusBadRequest)
		return
	}

	api.AddProducts(product)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Products", r.URL.Path)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert ID", http.StatusBadRequest)
		return
	}

	product := &api.Product{}

	err = product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusBadRequest)
		return
	}

	err = api.UpdateProducts(id, product)
	if err == api.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}
