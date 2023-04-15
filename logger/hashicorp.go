package logger

import (
	"github.com/hashicorp/go-hclog"
)

func New(name, level string) hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:       name,
		Level:      hclog.LevelFromString(level),
		Color:      1,
		TimeFormat: "02/01/2006 15:04:05",
	})
}
