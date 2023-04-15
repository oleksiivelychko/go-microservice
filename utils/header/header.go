package header

import "net/http"

func ContentTypeJSON(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
}
