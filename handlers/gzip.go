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

func NewGzipResponseWriter(responseWriter http.ResponseWriter) *GzipResponseWriter {
	gzipWriter := gzip.NewWriter(responseWriter)
	return &GzipResponseWriter{responseWriter: responseWriter, gzipWriter: gzipWriter}
}

func (gzipResponseWriter *GzipResponseWriter) Header() http.Header {
	return gzipResponseWriter.responseWriter.Header()
}

func (gzipResponseWriter *GzipResponseWriter) Write(bytes []byte) (int, error) {
	return gzipResponseWriter.gzipWriter.Write(bytes)
}

func (gzipResponseWriter *GzipResponseWriter) WriteHeader(statusCode int) {
	gzipResponseWriter.responseWriter.WriteHeader(statusCode)
}

func (gzipResponseWriter *GzipResponseWriter) Flush() {
	_ = gzipResponseWriter.gzipWriter.Flush()
	_ = gzipResponseWriter.gzipWriter.Close()
}

func (gzipHandler *GzipHandler) MiddlewareGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			gzipHandler.logger.Info("discovered `gzip` content-encoding")

			gzipResponseWriter := NewGzipResponseWriter(writer)
			gzipResponseWriter.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(gzipResponseWriter, request)
			defer gzipResponseWriter.Flush()

			return
		}

		next.ServeHTTP(writer, request)
	})
}
