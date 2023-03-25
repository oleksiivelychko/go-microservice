package product_handler

import (
	"context"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/errors"
	"github.com/oleksiivelychko/go-utils/io_json"
	"net/http"
)

func (productHandler *ProductHandler) MiddlewareValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")

		product := &api.Product{}

		err := io_json.FromJSON(product, request.Body)
		if err != nil {
			productHandler.logger.Error("deserializer", "error", err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			_ = io_json.ToJSON(&errors.GenericError{Message: err.Error()}, writer)
			return
		}

		validationErrors := productHandler.validation.Validate(product)
		if len(validationErrors) > 0 {
			productHandler.logger.Error("validation", "error", err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			// return the validation messages as an array
			_ = io_json.ToJSON(&errors.ValidationErrors{Messages: validationErrors.Errors()}, writer)
			return
		}

		// put product into context
		contextWith := context.WithValue(request.Context(), KeyProduct{}, product)
		request = request.WithContext(contextWith)

		// call the next handler, which can be another middleware in the chain, or the final handler
		next.ServeHTTP(writer, request)
	})
}

func (productHandler *ProductHandler) MiddlewareCurrency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		currency := request.URL.Query().Get("currency")
		if currency != "" {
			productHandler.productService.CurrencyService.SetCurrency(currency)
		}
		next.ServeHTTP(writer, request)
	})
}
