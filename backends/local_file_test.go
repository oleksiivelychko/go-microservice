package backends

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

const savePath = "/1/test.png"
const fileContent = "Hello World"

func setup(t *testing.T) (*Local, func()) {
	dir, err := os.MkdirTemp(os.TempDir(), "files")
	if err != nil {
		t.Fatal(err)
	}

	// 1mb = 1000000 bytes
	local, err := NewLocal(dir, 1000000)
	if err != nil {
		t.Fatal(err)
	}

	return local, func() {
		_ = os.RemoveAll(dir)
	}
}

func TestLocalSave(t *testing.T) {
	local, cleanup := setup(t)
	defer cleanup()

	content := bytes.NewBuffer([]byte(fileContent))
	writtenBytes, err := local.Save(savePath, content)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, writtenBytes, int64(unsafe.Sizeof(content)))

	file, err := os.Open(filepath.Join(local.basePath, savePath))
	assert.NoError(t, err)

	data, err := io.ReadAll(file)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, string(data))
}

func TestLocalGet(t *testing.T) {
	local, cleanup := setup(t)
	defer cleanup()

	content := bytes.NewBuffer([]byte(fileContent))
	_, err := local.Save(savePath, content)
	assert.NoError(t, err)

	file, err := local.Get(savePath)
	assert.NoError(t, err)
	defer file.Close()

	data, err := io.ReadAll(file)
	assert.Equal(t, fileContent, string(data))
}
