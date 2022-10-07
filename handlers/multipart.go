package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/contracts"
	"github.com/oleksiivelychko/go-microservice/utils"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

// MultipartHandler for creating and updating products as multipart/form-data
type MultipartHandler struct {
	log   hclog.Logger
	v     *utils.Validation
	store contracts.Storage
}

// NewMultipartHandler returns a new multipart handler with the given logger and validation
func NewMultipartHandler(l hclog.Logger, v *utils.Validation, s contracts.Storage) *MultipartHandler {
	return &MultipartHandler{l, v, s}
}

func (mp *MultipartHandler) ProcessForm(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024) // 32Mb
	if err != nil {
		mp.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		productId = api.GetNextProductId()
	}

	price, _ := strconv.ParseFloat(r.FormValue("price"), 32)
	product := api.Product{
		ID:          productId,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       float32(price),
		SKU:         r.FormValue("SKU"),
	}

	mpFile, mpHandler, err := r.FormFile("image")
	if err != nil {
		mp.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected file", http.StatusBadRequest)
		return
	}

	err = mp.saveFile(strconv.Itoa(productId), mpHandler.Filename, mpFile)
	if err != nil {
		mp.log.Error("unable to save file", "error", err)
		http.Error(rw, "unable to save file", http.StatusInternalServerError)
		return
	}

	if id == "" {
		api.AddProduct(product)
	} else {
		_ = api.UpdateProduct(product)
	}
}

func (mp *MultipartHandler) saveFile(id, path string, r io.ReadCloser) error {
	mp.log.Info("save file as part of multipart/form-data for the product", "id", id, "path", path)

	filePath := filepath.Join(id, path)
	_, err := mp.store.Save(filePath, r)

	return err
}
