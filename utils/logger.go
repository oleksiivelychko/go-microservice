package utils

import "github.com/hashicorp/go-hclog"

func NewLogger(name string) hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:       name,
		Level:      hclog.LevelFromString("debug"),
		Color:      1,
		TimeFormat: "15:04:05",
	})
}
