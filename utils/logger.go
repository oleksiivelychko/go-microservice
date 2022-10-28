package utils

import "github.com/hashicorp/go-hclog"

func NewLogger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:       "go-microservice",
		Level:      hclog.LevelFromString("debug"),
		Color:      1,
		TimeFormat: "15:04:05",
	})
}
