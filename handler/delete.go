package handler

import (
	"github.com/oleksiivelychko/go-microservice/utils"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product.
//
// responses:
// 204: noContentResponse
// 404: notFoundResponse
func (handler *Product) DeleteProduct(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Debug("DELETE /products")
	utils.HeaderContentTypeJSON(resp)

	id := handler.getProductID(req)

	productNotFoundErr := handler.productService.DeleteProduct(id)
	if productNotFoundErr != nil {
		handler.logger.Error(productNotFoundErr.Error())
		resp.WriteHeader(http.StatusNotFound)
		utils.ToJSON(&productNotFoundErr, resp)
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}
