package handlers

import (
	"github.com/oleksiivelychko/go-grpc-service/logger"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/services"
	"github.com/oleksiivelychko/go-microservice/storage"
	"github.com/oleksiivelychko/go-microservice/validation"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

const formDataMaxMemory32MB = 128 * 1024

// Multipart for CRUD actions regarding api.Product objects as multipart/form-data.
type Multipart struct {
	logger         *logger.Logger
	validation     *validation.Validate
	storage        storage.ILocal
	productService *services.Product
}

func NewMultipart(validation *validation.Validate, storage storage.ILocal, productService *services.Product, logger *logger.Logger) *Multipart {
	return &Multipart{logger, validation, storage, productService}
}

func (handler *Multipart) ProcessForm(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(formDataMaxMemory32MB)
	if err != nil {
		handler.logger.Error("expected multipart form data", "error", err)
		http.Error(resp, "expected multipart form data", http.StatusUnprocessableEntity)
		return
	}

	id := req.FormValue("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		productID = handler.productService.GetNextProductID()
	}

	price, err := strconv.ParseFloat(req.FormValue("price"), 64)
	if err != nil {
		handler.logger.Error("unable to parse price value", "error", err)
		http.Error(resp, "unable to parse price value", http.StatusUnprocessableEntity)
		return
	}

	product := api.Product{
		ID:    productID,
		Name:  req.FormValue("name"),
		Price: price,
		SKU:   req.FormValue("SKU"),
	}

	imageFile, fileHeader, err := req.FormFile("image")
	if err != nil {
		handler.logger.Error("expected image file", "error", err)
		http.Error(resp, "expected image file", http.StatusUnprocessableEntity)
		return
	}

	err = handler.saveFile(strconv.Itoa(productID), fileHeader.Filename, imageFile)
	if err != nil {
		handler.logger.Error("unable to save file", "error", err)
		http.Error(resp, "unable to save file", http.StatusInternalServerError)
		return
	}

	if id == "" {
		err = handler.productService.AddProduct(&product)
	} else {
		err = handler.productService.UpdateProduct(&product)
	}

	if err != nil {
		handler.logger.Error("req to gRPC service", "error", err)
		http.Error(resp, "req to gRPC service", http.StatusBadRequest)
	}
}

func (handler *Multipart) saveFile(id, path string, reader io.ReadCloser) error {
	filePath := filepath.Join(id, path)

	_, err := handler.storage.Save(filePath, reader)
	if err != nil {
		handler.logger.Info(
			"file from multipart/form-data has been successfully uploaded to",
			"filePath",
			filePath,
		)
	}

	return err
}
