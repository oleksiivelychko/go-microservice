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

func NewFileHandler(store contracts.Storage, log hclog.Logger) *File {
	return &File{storage: store, logger: log}
}

func (file *File) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// mux already has checked parameters according to regex rules
	vars := mux.Vars(r)
	productId := vars["id"]
	filename := vars["filename"]

	file.saveFile(productId, filename, rw, r)
}

func (file *File) invalidURI(uri string, rw http.ResponseWriter) {
	file.logger.Error("invalid path", "path", uri)
	http.Error(rw, "invalid file path: should be in the format: /[id]/[filename]", http.StatusBadRequest)
}

func (file *File) saveFile(id, filename string, rw http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(id, filename)
	_, err := file.storage.Save(filePath, r.Body)
	if err != nil {
		file.logger.Error("unable to save file", "error", err)
		http.Error(rw, "unable to save file", http.StatusInternalServerError)
	}

	file.logger.Info("file has been successfully uploaded to", "filePath", filePath)
}
