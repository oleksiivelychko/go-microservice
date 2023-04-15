package logger

import (
	"bytes"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLogger_Hashicorp(t *testing.T) {
	var buffer bytes.Buffer

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.LevelFromString("DEBUG"),
		Output: &buffer,
	})

	logger.Debug("this is test", "who", "programmer", "why", "testing")

	bufStr := buffer.String()
	idxByte := strings.IndexByte(bufStr, ' ')
	rest := bufStr[idxByte+1:]

	assert.Equal(t, "[DEBUG] test: this is test: who=programmer why=testing\n", rest)
}
