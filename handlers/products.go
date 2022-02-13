package handlers

import (
	"github.com/oleksiivelychko/go-microservice/api"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProducts(w, r)
	case http.MethodPost:
		p.addProduct(w, r)
	case http.MethodPut:
		p.updateProduct(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products", r.URL.Path)

	lp := api.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products", r.URL.Path)

	product := &api.Product{}

	err := product.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusBadRequest)
		return
	}

	api.AddProducts(product)
}

func (p *Products) updateProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Products", r.URL.Path)

	// expects ID in URL
	regex := regexp.MustCompile("/([0-9]+)")
	g := regex.FindAllStringSubmatch(r.URL.Path, -1)
	if len(g) != 1 {
		p.l.Println("Invalid URI has more than one ID")
		http.Error(w, "Invalid URI", http.StatusBadRequest)
		return
	}

	if len(g[0]) != 2 {
		p.l.Println("Invalid URI has more than one capture group")
		http.Error(w, "Invalid URI", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(g[0][1])
	if err != nil {
		p.l.Println("Invalid URI has id that unable to convert to integer")
		http.Error(w, "Invalid URI", http.StatusBadRequest)
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
