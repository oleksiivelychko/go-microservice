package handlers

import (
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/service"
	"github.com/oleksiivelychko/go-microservice/utils"
	"github.com/oleksiivelychko/go-utils/storage"
	"github.com/oleksiivelychko/go-utils/validation"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

// MultipartHandler for CRUD actions regarding api.Product objects as multipart/form-data.
type MultipartHandler struct {
	logger         hclog.Logger
	validation     *validation.Validate
	storage        storage.ILocal
	productService *service.ProductService
}

func NewMultipartHandler(
	logger hclog.Logger,
	validation *validation.Validate,
	storage storage.ILocal,
	productService *service.ProductService,
) *MultipartHandler {
	return &MultipartHandler{logger, validation, storage, productService}
}

func (handler *MultipartHandler) ProcessForm(responseWriter http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(utils.FormDataMaxMemory32MB)
	if err != nil {
		handler.logger.Error("expected multipart form data", "error", err)
		http.Error(responseWriter, "expected multipart form data", http.StatusUnprocessableEntity)
		return
	}

	id := request.FormValue("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		productID = handler.productService.GetNextProductID()
	}

	price, err := strconv.ParseFloat(request.FormValue("price"), 64)
	if err != nil {
		handler.logger.Error("unable to parse price value", "error", err)
		http.Error(responseWriter, "unable to parse price value", http.StatusUnprocessableEntity)
		return
	}

	product := api.Product{
		ID:    productID,
		Name:  request.FormValue("name"),
		Price: price,
		SKU:   request.FormValue("SKU"),
	}

	imageFile, fileHeader, err := request.FormFile("image")
	if err != nil {
		handler.logger.Error("expected image file", "error", err)
		http.Error(responseWriter, "expected image file", http.StatusUnprocessableEntity)
		return
	}

	err = handler.saveFile(strconv.Itoa(productID), fileHeader.Filename, imageFile)
	if err != nil {
		handler.logger.Error("unable to save file", "error", err)
		http.Error(responseWriter, "unable to save file", http.StatusInternalServerError)
		return
	}

	if id == "" {
		err = handler.productService.AddProduct(&product)
	} else {
		err = handler.productService.UpdateProduct(&product)
	}

	if err != nil {
		handler.logger.Error("request to gRPC service", "error", err)
		http.Error(responseWriter, "request to gRPC service", http.StatusBadRequest)
	}
}

func (handler *MultipartHandler) saveFile(id, path string, readCloser io.ReadCloser) error {
	filePath := filepath.Join(id, path)

	_, err := handler.storage.Save(filePath, readCloser)
	if err != nil {
		handler.logger.Info(
			"file from multipart/form-data has been successfully uploaded to",
			"filePath",
			filePath,
		)
	}

	return err
}
