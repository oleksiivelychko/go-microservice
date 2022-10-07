package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-microservice/contracts"
	"net/http"
	"path/filepath"
)

type File struct {
	log   hclog.Logger
	store contracts.Storage
}

func NewFile(store contracts.Storage, log hclog.Logger) *File {
	return &File{store: store, log: log}
}

func (file *File) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// mux already has checked parameters according to regex rules
	vars := mux.Vars(r)
	productId := vars["id"]
	filename := vars["filename"]

	file.log.Info("handle POST", "id", productId, "filename", filename)
	file.saveFile(productId, filename, rw, r)
}

func (file *File) UploadMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32Mb
	if err != nil {
		file.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	file.log.Info("Process multipart/form data for ID", "ID", id)
}

func (file *File) invalidURI(uri string, rw http.ResponseWriter) {
	file.log.Error("invalid path", "path", uri)
	http.Error(rw, "invalid file path: should be in the format: /[id]/[filename]", http.StatusBadRequest)
}

func (file *File) saveFile(id, path string, rw http.ResponseWriter, r *http.Request) {
	file.log.Info("save file for the product", "id", id, "path", path)

	filePath := filepath.Join(id, path)
	_, err := file.store.Save(filePath, r.Body)
	if err != nil {
		file.log.Error("unable to save file", "error", err)
		http.Error(rw, "unable to save file", http.StatusInternalServerError)
	}
}
