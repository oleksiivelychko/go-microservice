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
	storage storage.ILocal
}

func NewFile(storage storage.ILocal, logger hclog.Logger) *File {
	return &File{storage: storage, logger: logger}
}

func (handler *File) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// mux already has checked parameters according to regex rules
	muxVars := mux.Vars(req)
	filePath := filepath.Join(muxVars["id"], muxVars["filename"])

	_, err := handler.storage.Save(filePath, req.Body)
	if err != nil {
		handler.logger.Error("unable to save file", "error", err)
		http.Error(resp, "unable to save file", http.StatusInternalServerError)
	}

	handler.logger.Info("file has been successfully uploaded to", "filePath", filePath)
}
