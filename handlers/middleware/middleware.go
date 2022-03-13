package middleware

import (
	"context"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-microservice/handlers"
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

/*
ProductHandlerEmbed Cannot extend types defined in other packages.
Solution is embed a type in another package in your own type, then extend your own type.
*/
type ProductHandlerEmbed struct {
	handlers.ProductHandler
}

func (p *ProductHandlerEmbed) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := &api.Product{}

		err := utils.FromJSON(product, r.Body)
		if err != nil {
			p.L.Println("[ERROR] deserializing product", err)

			rw.WriteHeader(http.StatusBadRequest)
			_ = utils.ToJSON(&handlers.GenericError{Message: err.Error()}, rw)

			return
		}

		errs := p.V.Validate(product)
		if len(errs) != 0 {
			p.L.Println("[ERROR] validating product", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)

			_ = utils.ToJSON(&handlers.ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), handlers.KeyProduct{}, product)
		r = r.WithContext(ctx)

		// call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
