package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/contracts"
	"github.com/oleksiivelychko/go-microservice/service"
	validatorHelper "github.com/oleksiivelychko/go-utils/validator_helper"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

// MultipartHandler for CRUD actions regarding api.Product objects as multipart/form-data.
type MultipartHandler struct {
	logger         hclog.Logger
	validation     *validatorHelper.Validation
	storage        contracts.Storage
	productService *service.ProductService
}

func NewMultipartHandler(l hclog.Logger, v *validatorHelper.Validation, s contracts.Storage, ps *service.ProductService) *MultipartHandler {
	return &MultipartHandler{l, v, s, ps}
}

func (handler *MultipartHandler) ProcessForm(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(128 * 1024) // 32Mb
	if err != nil {
		handler.logger.Error("expected multipart form data", "error", err)
		http.Error(writer, "expected multipart form data", http.StatusUnprocessableEntity)
		return
	}

	id := request.FormValue("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		productId = handler.productService.GetNextProductId()
	}

	price, err := strconv.ParseFloat(request.FormValue("price"), 64)
	if err != nil {
		handler.logger.Error("unable to parse price value to float type", "error", err)
		http.Error(writer, "unable to parse price value to float type", http.StatusUnprocessableEntity)
		return
	}

	product := api.Product{
		ID:    productId,
		Name:  request.FormValue("name"),
		Price: price,
		SKU:   request.FormValue("SKU"),
	}

	imageFile, fileHeader, err := request.FormFile("image")
	if err != nil {
		handler.logger.Error("expected file", "error", err)
		http.Error(writer, "expected file", http.StatusUnprocessableEntity)
		return
	}

	err = handler.saveFile(strconv.Itoa(productId), fileHeader.Filename, imageFile)
	if err != nil {
		handler.logger.Error("unable to save file", "error", err)
		http.Error(writer, "unable to save file", http.StatusInternalServerError)
		return
	}

	if id == "" {
		err = handler.productService.AddProduct(&product)
	} else {
		err = handler.productService.UpdateProduct(&product)
	}

	if err != nil {
		handler.logger.Error("request to gRPC service", "error", err)
		http.Error(writer, "request to gRPC service", http.StatusBadRequest)
	}
}

func (handler *MultipartHandler) saveFile(id, path string, readCloser io.ReadCloser) error {
	filePath := filepath.Join(id, path)

	_, err := handler.storage.Save(filePath, readCloser)
	if err != nil {
		handler.logger.Info("file from multipart/form-data has been successfully uploaded to", "filePath", filePath)
	}

	return err
}
