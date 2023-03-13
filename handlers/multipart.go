package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-utils/storage"
	validatorUtils "github.com/oleksiivelychko/go-utils/validator_helper"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

// MultipartHandler for CRUD actions regarding api.Product objects as multipart/form-data.
type MultipartHandler struct {
	logger         hclog.Logger
	validation     *validatorUtils.Validation
	storage        storage.LocalStorage
	productService *service.ProductService
}

func NewMultipartHandler(
	logger hclog.Logger,
	validation *validatorUtils.Validation,
	storage storage.LocalStorage,
	productService *service.ProductService,
) *MultipartHandler {
	return &MultipartHandler{logger, validation, storage, productService}
}

func (multipartHandler *MultipartHandler) ProcessForm(responseWriter http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(128 * 1024) // 32Mb
	if err != nil {
		multipartHandler.logger.Error("expected multipart form data", "error", err)
		http.Error(responseWriter, "expected multipart form data", http.StatusUnprocessableEntity)
		return
	}

	id := request.FormValue("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		productId = multipartHandler.productService.GetNextProductId()
	}

	price, err := strconv.ParseFloat(request.FormValue("price"), 64)
	if err != nil {
		multipartHandler.logger.Error("unable to parse price value to float type", "error", err)
		http.Error(responseWriter, "unable to parse price value to float type", http.StatusUnprocessableEntity)
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
		multipartHandler.logger.Error("expected file", "error", err)
		http.Error(responseWriter, "expected file", http.StatusUnprocessableEntity)
		return
	}

	err = multipartHandler.saveFile(strconv.Itoa(productId), fileHeader.Filename, imageFile)
	if err != nil {
		multipartHandler.logger.Error("unable to save file", "error", err)
		http.Error(responseWriter, "unable to save file", http.StatusInternalServerError)
		return
	}

	if id == "" {
		err = multipartHandler.productService.AddProduct(&product)
	} else {
		err = multipartHandler.productService.UpdateProduct(&product)
	}

	if err != nil {
		multipartHandler.logger.Error("request to gRPC service", "error", err)
		http.Error(responseWriter, "request to gRPC service", http.StatusBadRequest)
	}
}

func (multipartHandler *MultipartHandler) saveFile(id, path string, readCloser io.ReadCloser) error {
	filePath := filepath.Join(id, path)

	_, err := multipartHandler.storage.Save(filePath, readCloser)
	if err != nil {
		multipartHandler.logger.Info("file from multipart/form-data has been successfully uploaded to", "filePath", filePath)
	}

	return err
}
