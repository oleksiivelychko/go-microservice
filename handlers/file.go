package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-utils/storage"
	"net/http"
	"path/filepath"
)

type File struct {
	logger  hclog.Logger
	storage storage.LocalStorage
}

func NewFileHandler(storage storage.LocalStorage, logger hclog.Logger) *File {
	return &File{storage: storage, logger: logger}
}

func (file *File) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// mux already has checked parameters according to regex rules
	muxVars := mux.Vars(request)
	productId := muxVars["id"]
	filename := muxVars["filename"]

	file.saveFile(productId, filename, responseWriter, request)
}

func (file *File) invalidURI(uri string, responseWriter http.ResponseWriter) {
	file.logger.Error("invalid path", "path", uri)
	http.Error(responseWriter, "invalid file path: should be in the format: /[id]/[filename]", http.StatusBadRequest)
}

func (file *File) saveFile(id, filename string, responseWriter http.ResponseWriter, request *http.Request) {
	filePath := filepath.Join(id, filename)
	_, err := file.storage.Save(filePath, request.Body)
	if err != nil {
		file.logger.Error("unable to save file", "error", err)
		http.Error(responseWriter, "unable to save file", http.StatusInternalServerError)
	}

	file.logger.Info("file has been successfully uploaded to", "filePath", filePath)
}
