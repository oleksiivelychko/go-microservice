package handler

import (
	"context"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/errors"
	"github.com/oleksiivelychko/go-microservice/utils"
	"github.com/oleksiivelychko/go-utils/response"
	"github.com/oleksiivelychko/go-utils/serializer"
	"net/http"
)

func (handler *ProductHandler) MiddlewareValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		response.HeaderContentTypeJSON(responseWriter)

		product := &api.Product{}

		err := serializer.FromJSON(product, request.Body)
		if err != nil {
			handler.logger.Error("deserializer", "error", err)
			responseWriter.WriteHeader(http.StatusUnprocessableEntity)
			_ = serializer.ToJSON(&errors.GenericError{Message: err.Error()}, responseWriter)
			return
		}

		validationErrors := handler.validation.Validate(product)
		if len(validationErrors) > 0 {
			handler.logger.Error("validation", "error", err)
			responseWriter.WriteHeader(http.StatusUnprocessableEntity)
			// return the validation messages as an array
			_ = serializer.ToJSON(&errors.ValidationErrors{Messages: validationErrors.Errors()}, responseWriter)
			return
		}

		// put product into context
		contextWith := context.WithValue(request.Context(), KeyProduct{}, product)
		request = request.WithContext(contextWith)

		// call the next handler, which can be another middleware in the chain, or the final handler
		next.ServeHTTP(responseWriter, request)
	})
}

func (handler *ProductHandler) MiddlewareCurrency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		currency := request.URL.Query().Get(utils.CurrencyQueryParam)
		if currency != "" {
			handler.productService.CurrencyService.SetCurrency(currency)
		}
		next.ServeHTTP(writer, request)
	})
}
