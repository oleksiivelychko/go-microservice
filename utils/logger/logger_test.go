package logger

import (
	"bytes"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLogger_Hashicorp(t *testing.T) {
	var buf bytes.Buffer

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.LevelFromString("DEBUG"),
		Output: &buf,
	})

	logger.Debug("this is test", "who", "programmer", "why", "testing")

	bufStr := buf.String()
	idxByte := strings.IndexByte(bufStr, ' ')
	rest := bufStr[idxByte+1:]

	assert.Equal(t, "[DEBUG] test: this is test: who=programmer why=testing\n", rest)
}
