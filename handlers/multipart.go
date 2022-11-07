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

// MultipartHandler for CRUD actions regarding api.Product objects as multipart/form-data.
type MultipartHandler struct {
	log hclog.Logger
	val *utils.Validation
	stg contracts.Storage
	srv *service.ProductService
}

func NewMultipartHandler(l hclog.Logger, v *utils.Validation, s contracts.Storage, ps *service.ProductService) *MultipartHandler {
	return &MultipartHandler{l, v, s, ps}
}

func (mp *MultipartHandler) ProcessForm(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024) // 32Mb
	if err != nil {
		mp.log.Error("expected multipart form data", "error", err)
		http.Error(rw, "expected multipart form data", http.StatusUnprocessableEntity)
		return
	}

	id := r.FormValue("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		productId = mp.srv.GetNextProductId()
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		mp.log.Error("unable to parse price value to float type", "error", err)
		http.Error(rw, "unable to parse price value to float type", http.StatusUnprocessableEntity)
		return
	}

	product := api.Product{
		ID:    productId,
		Name:  r.FormValue("name"),
		Price: price,
		SKU:   r.FormValue("SKU"),
	}

	imageFile, fileHeader, err := r.FormFile("image")
	if err != nil {
		mp.log.Error("expected file", "error", err)
		http.Error(rw, "expected file", http.StatusUnprocessableEntity)
		return
	}

	err = mp.saveFile(strconv.Itoa(productId), fileHeader.Filename, imageFile)
	if err != nil {
		mp.log.Error("unable to save file", "error", err)
		http.Error(rw, "unable to save file", http.StatusInternalServerError)
		return
	}

	if id == "" {
		err = mp.srv.AddProduct(&product)
	} else {
		err = mp.srv.UpdateProduct(&product)
	}

	if err != nil {
		mp.log.Error("request to gRPC service", "error", err)
		http.Error(rw, "request to gRPC service", http.StatusBadRequest)
	}
}

func (mp *MultipartHandler) saveFile(id, path string, r io.ReadCloser) error {
	filePath := filepath.Join(id, path)

	_, err := mp.stg.Save(filePath, r)
	if err != nil {
		mp.log.Info("file from multipart/form-data has been successfully uploaded to", "filePath", filePath)
	}

	return err
}
