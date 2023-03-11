package handlers

import (
	"context"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

func (handler *ProductHandler) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		writer.Header().Add("Content-Type", "application/json")

		product := &api.Product{}

		err := utils.FromJSON(product, r.Body)
		if err != nil {
			handler.logger.Error("JSON decode", "error", err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			_ = utils.ToJSON(&GenericError{Message: err.Error()}, writer)
			return
		}

		errs := handler.validation.Validate(product)
		if len(errs) != 0 {
			handler.logger.Error("validation", "error", err)
			writer.WriteHeader(http.StatusUnprocessableEntity)
			// return the validation messages as an array
			_ = utils.ToJSON(&ValidationErrors{Messages: errs.Errors()}, writer)
			return
		}

		// add the product into the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)

		// call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(writer, r)
	})
}

func (handler *ProductHandler) MiddlewareProductCurrency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		currency := request.URL.Query().Get("currency")
		if currency != "" {
			handler.productService.Currency.SetCurrency(currency)
		}
		next.ServeHTTP(writer, request)
	})
}
