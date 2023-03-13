package handlers

import (
	"context"
	"github.com/oleksiivelychko/go-microservice/api"
	jsonUtils "github.com/oleksiivelychko/go-utils/json_io"
	"net/http"
)

func (productHandler *ProductHandler) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")

		product := &api.Product{}

		err := jsonUtils.FromJSON(product, request.Body)
		if err != nil {
			productHandler.logger.Error("JSON decode", "error", err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			_ = jsonUtils.ToJSON(&GenericError{Message: err.Error()}, writer)
			return
		}

		validationErrors := productHandler.validation.Validate(product)
		if len(validationErrors) != 0 {
			productHandler.logger.Error("validation", "error", err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			// return the validation messages as an array
			_ = jsonUtils.ToJSON(&ValidationErrors{Messages: validationErrors.Errors()}, writer)
			return
		}

		// add the product into the context
		ctx := context.WithValue(request.Context(), KeyProduct{}, product)
		request = request.WithContext(ctx)

		// call the next productHandler, which can be another middleware in the chain, or the final productHandler.
		next.ServeHTTP(writer, request)
	})
}

func (productHandler *ProductHandler) MiddlewareProductCurrency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		currency := request.URL.Query().Get("currency")
		if currency != "" {
			productHandler.productService.CurrencyService.SetCurrency(currency)
		}
		next.ServeHTTP(writer, request)
	})
}
