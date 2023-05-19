package handlers

import (
	"compress/gzip"
	"github.com/oleksiivelychko/go-grpc-service/logger"
	"net/http"
	"strings"
)

type GZIP struct {
	logger *logger.Logger
}

func NewGZIP(logger *logger.Logger) *GZIP {
	return &GZIP{logger}
}

type ResponseGZIP struct {
	responseWriter http.ResponseWriter
	gzipWriter     *gzip.Writer
}

func NewResponseGZIP(resp http.ResponseWriter) *ResponseGZIP {
	return &ResponseGZIP{responseWriter: resp, gzipWriter: gzip.NewWriter(resp)}
}

func (resp *ResponseGZIP) Header() http.Header {
	return resp.responseWriter.Header()
}

func (resp *ResponseGZIP) Write(bytes []byte) (int, error) {
	return resp.gzipWriter.Write(bytes)
}

func (resp *ResponseGZIP) WriteHeader(statusCode int) {
	resp.responseWriter.WriteHeader(statusCode)
}

func (resp *ResponseGZIP) Flush() {
	resp.gzipWriter.Flush()
	resp.gzipWriter.Close()
}

func (handler *GZIP) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			handler.logger.Info("discovered gzip content-encoding")

			respGZIP := NewResponseGZIP(resp)
			resp.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(respGZIP, req)
			defer respGZIP.Flush()

			return
		}

		next.ServeHTTP(resp, req)
	})
}
