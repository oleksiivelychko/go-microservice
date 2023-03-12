package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/contracts"
	"net/http"
	"path/filepath"
)

type File struct {
	logger  hclog.Logger
	storage contracts.Storage
}

func NewFileHandler(storage contracts.Storage, logger hclog.Logger) *File {
	return &File{storage: storage, logger: logger}
}

func (file *File) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// mux already has checked parameters according to regex rules
	vars := mux.Vars(request)
	productId := vars["id"]
	filename := vars["filename"]

	file.saveFile(productId, filename, writer, request)
}

func (file *File) invalidURI(uri string, writer http.ResponseWriter) {
	file.logger.Error("invalid path", "path", uri)
	http.Error(writer, "invalid file path: should be in the format: /[id]/[filename]", http.StatusBadRequest)
}

func (file *File) saveFile(id, filename string, writer http.ResponseWriter, request *http.Request) {
	filePath := filepath.Join(id, filename)
	_, err := file.storage.Save(filePath, request.Body)
	if err != nil {
		file.logger.Error("unable to save file", "error", err)
		http.Error(writer, "unable to save file", http.StatusInternalServerError)
	}

	file.logger.Info("file has been successfully uploaded to", "filePath", filePath)
}
