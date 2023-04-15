package utils

import "net/http"

func HeaderContentTypeJSON(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
}
