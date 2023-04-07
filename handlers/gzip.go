package handlers

import (
	"compress/gzip"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-utils/response"
	"net/http"
	"strings"
)

type HandlerGZIP struct {
	logger hclog.Logger
}

func NewHandlerGZIP(logger hclog.Logger) *HandlerGZIP {
	return &HandlerGZIP{logger}
}

type ResponseWriterGZIP struct {
	responseWriter http.ResponseWriter
	gzipWriter     *gzip.Writer
}

func NewResponseWriterGZIP(responseWriter http.ResponseWriter) *ResponseWriterGZIP {
	return &ResponseWriterGZIP{responseWriter: responseWriter, gzipWriter: gzip.NewWriter(responseWriter)}
}

func (gzipResponseWriter *ResponseWriterGZIP) Header() http.Header {
	return gzipResponseWriter.responseWriter.Header()
}

func (gzipResponseWriter *ResponseWriterGZIP) Write(bytes []byte) (int, error) {
	return gzipResponseWriter.gzipWriter.Write(bytes)
}

func (gzipResponseWriter *ResponseWriterGZIP) WriteHeader(statusCode int) {
	gzipResponseWriter.responseWriter.WriteHeader(statusCode)
}

func (gzipResponseWriter *ResponseWriterGZIP) Flush() {
	_ = gzipResponseWriter.gzipWriter.Flush()
	_ = gzipResponseWriter.gzipWriter.Close()
}

func (handler *HandlerGZIP) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			handler.logger.Info("discovered gzip content-encoding")

			gzipResponseWriter := NewResponseWriterGZIP(writer)
			response.HeaderContentEncodingGZIP(gzipResponseWriter)

			next.ServeHTTP(gzipResponseWriter, request)
			defer gzipResponseWriter.Flush()

			return
		}

		next.ServeHTTP(writer, request)
	})
}
