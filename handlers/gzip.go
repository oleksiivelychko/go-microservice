package handlers

import (
	"compress/gzip"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strings"
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

func (g *GzipHandler) MiddlewareGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			g.log.Info("discovered `gzip` content-encoding")

			wrw := NewGzipResponseWriter(rw)
			wrw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(wrw, r)
			defer wrw.Flush()

			return
		}

		next.ServeHTTP(rw, r)
	})
}
