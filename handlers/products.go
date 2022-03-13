package handlers

import (
	"github.com/oleksiivelychko/go-microservice/helpers"
	"log"
)

type Products struct {
	l *log.Logger
	v *helpers.Validation
}

type KeyProduct struct{}
