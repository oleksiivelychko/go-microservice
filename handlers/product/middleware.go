package product

import (
	"context"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/errors"
	"github.com/oleksiivelychko/go-microservice/utils/header"
	"github.com/oleksiivelychko/go-microservice/utils/serializer"
	"net/http"
)

func (handler *Handler) MiddlewareValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		header.ContentTypeJSON(resp)

		product := &api.Product{}

		err := serializer.FromJSON(product, req.Body)
		if err != nil {
			handler.logger.Error("deserializer", "error", err)
			resp.WriteHeader(http.StatusUnprocessableEntity)
			serializer.ToJSON(&errors.GenericError{Message: err.Error()}, resp)
			return
		}

		validationErrors := handler.validation.Validate(product)
		if len(validationErrors) > 0 {
			handler.logger.Error("validation", "error", err)
			resp.WriteHeader(http.StatusUnprocessableEntity)
			// return the validation messages as an array
			serializer.ToJSON(&errors.ValidationErrors{Messages: validationErrors.Errors()}, resp)
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
