package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/contracts"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-microservice/utils"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

// MultipartHandler for creating and updating products as multipart/form-data.
type MultipartHandler struct {
	l     hclog.Logger
	v     *utils.Validation
	store contracts.Storage
	ps    *service.ProductService
}

// NewMultipartHandler returns a new multipart handler with the given logger and validation.
func NewMultipartHandler(
	l hclog.Logger,
	v *utils.Validation,
	s contracts.Storage,
	ps *service.ProductService,
) *MultipartHandler {
	return &MultipartHandler{l, v, s, ps}
}

func (mp *MultipartHandler) ProcessForm(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024) // 32Mb
	if err != nil {
		mp.l.Error("expected multipart form data", "error", err)
		http.Error(rw, "expected multipart form data", http.StatusUnprocessableEntity)
		return
	}

	id := r.FormValue("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		productId = mp.ps.GetNextProductId()
	}

	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	product := api.Product{
		ID:    productId,
		Name:  r.FormValue("name"),
		Price: price,
		SKU:   r.FormValue("SKU"),
	}

	mpFile, mpHandler, err := r.FormFile("image")
	if err != nil {
		mp.l.Error("expected file", "error", err)
		http.Error(rw, "expected file", http.StatusUnprocessableEntity)
		return
	}

	err = mp.saveFile(strconv.Itoa(productId), mpHandler.Filename, mpFile)
	if err != nil {
		mp.l.Error("unable to save file", "error", err)
		http.Error(rw, "unable to save file", http.StatusInternalServerError)
		return
	}

	if id == "" {
		err = mp.ps.AddProduct(&product)
	} else {
		err = mp.ps.UpdateProduct(&product)
	}

	if err != nil {
		mp.l.Error("unable to make request to gRPC service", "error", err)
		http.Error(rw, "unable to make request to gRPC service", http.StatusBadRequest)
	}
}

func (mp *MultipartHandler) saveFile(id, path string, r io.ReadCloser) error {
	mp.l.Info("save file from multipart/form-data", "productId", id, "filePath", path)

	filePath := filepath.Join(id, path)
	_, err := mp.store.Save(filePath, r)

	return err
}
