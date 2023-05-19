package product

import (
	"context"
	"encoding/json"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/errors"
	"net/http"
)

func (handler *Handler) MiddlewareValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")

		product := &api.Product{}

		err := json.NewDecoder(req.Body).Decode(product)
		if err != nil {
			handler.logger.Error("deserializer", "error", err)
			resp.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(resp).Encode(&errors.GenericError{Message: err.Error()})
			return
		}

		validationErrors := handler.validation.Validate(product)
		if len(validationErrors) > 0 {
			handler.logger.Error("validation", "error", err)
			resp.WriteHeader(http.StatusUnprocessableEntity)
			// return the validation messages as an array
			json.NewEncoder(resp).Encode(&errors.ValidationErrors{Messages: validationErrors.Errors()})
			return
		}

		// put product into context
		contextWith := context.WithValue(req.Context(), KeyProduct{}, product)
		req = req.WithContext(contextWith)

		// call the next handler, which can be another middleware in the chain, or the final handler
		next.ServeHTTP(resp, req)
	})
}

func (handler *Handler) MiddlewareCurrency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		currency := req.URL.Query().Get("currency")
		if currency != "" {
			handler.productService.CurrencyService.SetCurrency(currency)
		}
		next.ServeHTTP(resp, req)
	})
}
