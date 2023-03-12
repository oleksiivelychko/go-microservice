package handlers

import (
	"compress/gzip"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strings"
)

type GzipHandler struct {
	logger hclog.Logger
}

func NewGzipHandler(logger hclog.Logger) *GzipHandler {
	return &GzipHandler{logger}
}

type GzipResponseWriter struct {
	responseWriter http.ResponseWriter
	gzipWriter     *gzip.Writer
}

func NewGzipResponseWriter(writer http.ResponseWriter) *GzipResponseWriter {
	gzipWriter := gzip.NewWriter(writer)
	return &GzipResponseWriter{responseWriter: writer, gzipWriter: gzipWriter}
}

func (writer *GzipResponseWriter) Header() http.Header {
	return writer.responseWriter.Header()
}

func (writer *GzipResponseWriter) Write(bytes []byte) (int, error) {
	return writer.gzipWriter.Write(bytes)
}

func (writer *GzipResponseWriter) WriteHeader(statuscode int) {
	writer.responseWriter.WriteHeader(statuscode)
}

func (writer *GzipResponseWriter) Flush() {
	_ = writer.gzipWriter.Flush()
	_ = writer.gzipWriter.Close()
}

func (handler *GzipHandler) MiddlewareGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			handler.logger.Info("discovered `gzip` content-encoding")

			gzipResponseWriter := NewGzipResponseWriter(writer)
			gzipResponseWriter.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(gzipResponseWriter, request)
			defer gzipResponseWriter.Flush()

			return
		}

		next.ServeHTTP(writer, request)
	})
}
