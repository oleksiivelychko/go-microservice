package handlers

import (
	"context"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
	"strings"
)

func (ph *ProductHandler) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		product := &api.Product{}

		err := utils.FromJSON(product, r.Body)
		if err != nil {
			ph.l.Error("deserialization", "error", err)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			_ = utils.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		errs := ph.v.Validate(product)
		if len(errs) != 0 {
			ph.l.Error("validation", "error", err)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			// return the validation messages as an array
			_ = utils.ToJSON(&ValidationErrors{Messages: errs.Errors()}, rw)
			return
		}

		// add the product into the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
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
