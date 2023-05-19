package server

import (
	"context"
	"net/http"
	"testing"
	"time"
)

type Handler struct{}

func (handler *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello, World!"))
}

func TestServer_ListenAndServeHTTP(t *testing.T) {
	server := NewHTTP("localhost:8080", &Handler{}, nil)

	go func() {
		time.Sleep(1 * time.Second)
		_ = server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		t.Fatal(err.Error())
	}
}
