package server

import (
	"github.com/hashicorp/go-hclog"
	"net/http"
	"time"
)

func NewHTTP(addr string, handler http.Handler, errorLog hclog.Logger) *http.Server {
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
		ReadTimeout:  10 * time.Second,  // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		ErrorLog:     errorLog.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true}),
	}

	return server
}
