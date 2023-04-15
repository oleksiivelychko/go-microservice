package header

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeader_ContentTypeJSON(t *testing.T) {
	handler := createHandler()

	req := httptest.NewRequest("GET", "https://example.com", nil)
	resp := httptest.NewRecorder()

	ContentTypeJSON(resp)
	handler(resp, req)

	response := resp.Result()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("unable to get successful status code, got %d", response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "application/json" {
		t.Errorf("content type mismatch: %s", response.Header.Get("Content-Type"))
	}

	if len(string(body)) == 0 {
		t.Error("unable to fetch data")
	}
}

func createHandler() func(resp http.ResponseWriter, req *http.Request) {
	handler := func(resp http.ResponseWriter, req *http.Request) {
		io.WriteString(resp, "Hello, World!")
	}

	return handler
}
