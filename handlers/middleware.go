package handlers

import (
	"context"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

func (ph *ProductHandler) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		product := &api.Product{}

		err := utils.FromJSON(product, r.Body)
		if err != nil {
			ph.log.Error("JSON decode", "error", err)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			_ = utils.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		errs := ph.val.Validate(product)
		if len(errs) != 0 {
			ph.log.Error("validation", "error", err)
			rw.WriteHeader(http.StatusUnprocessableEntity)
			// return the validation messages as an array
			_ = utils.ToJSON(&ValidationErrors{Messages: errs.Errors()}, rw)
			return
		}

		// add the product into the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// hand over currency into service
		ph.setCurrency(r)

		// call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}

func (ph *ProductHandler) MiddlewareProductCurrency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// hand over currency into service
		ph.setCurrency(r)
		next.ServeHTTP(rw, r)
	})
}
