package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-utils/storage"
	"net/http"
	"path/filepath"
)

type FileHandler struct {
	logger  hclog.Logger
	storage storage.ILocal
}

func NewFileHandler(storage storage.ILocal, logger hclog.Logger) *FileHandler {
	return &FileHandler{storage: storage, logger: logger}
}

func (fileHandler *FileHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// mux already has checked parameters according to regex rules
	muxVars := mux.Vars(request)
	filePath := filepath.Join(muxVars["id"], muxVars["filename"])

	_, err := fileHandler.storage.Save(filePath, request.Body)
	if err != nil {
		fileHandler.logger.Error("unable to save file", "error", err)
		http.Error(responseWriter, "unable to save file", http.StatusInternalServerError)
	}

	fileHandler.logger.Info("file has been successfully uploaded to", "filePath", filePath)
}
