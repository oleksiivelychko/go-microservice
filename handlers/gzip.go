package handlers

import (
	"compress/gzip"
	"github.com/hashicorp/go-hclog"
	"net/http"
)

type GzipHandler struct {
	log hclog.Logger
}

func NewGzipHandler(l hclog.Logger) *GzipHandler {
	return &GzipHandler{l}
}

type GzipResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewGzipResponseWriter(rw http.ResponseWriter) *GzipResponseWriter {
	gw := gzip.NewWriter(rw)
	return &GzipResponseWriter{rw: rw, gw: gw}
}

func (grw *GzipResponseWriter) Header() http.Header {
	return grw.rw.Header()
}

func (grw *GzipResponseWriter) Write(d []byte) (int, error) {
	return grw.gw.Write(d)
}

func (grw *GzipResponseWriter) WriteHeader(statuscode int) {
	grw.rw.WriteHeader(statuscode)
}

func (grw *GzipResponseWriter) Flush() {
	_ = grw.gw.Flush()
	_ = grw.gw.Close()
}
