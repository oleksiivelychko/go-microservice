package product

import (
	"encoding/json"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a product.
//
// responses:
// 204: noContentResponse
// 404: notFoundResponse
func (handler *Handler) DeleteProduct(resp http.ResponseWriter, req *http.Request) {
	handler.logger.Info("DELETE /products")
	resp.Header().Set("Content-Type", "application/json")

	id := handler.getProductID(req)
	productNotFoundErr := handler.productService.DeleteProduct(id)
	if productNotFoundErr != nil {
		handler.logger.Error(productNotFoundErr.Error())
		resp.WriteHeader(http.StatusNotFound)
		json.NewEncoder(resp).Encode(&productNotFoundErr)
		return
	}

	resp.WriteHeader(http.StatusNoContent)
}
