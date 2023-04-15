package storage

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

const savePath = "test.png"
const fileContent = "Hello, World!"

func setup(t *testing.T) (*Local, func()) {
	tmpDirectory, err := os.MkdirTemp(os.TempDir(), "files")
	if err != nil {
		t.Fatal(err)
	}

	// 1MB = 1000000 bytes
	storage, err := New(tmpDirectory, 1000000)
	if err != nil {
		t.Fatal(err)
	}

	return storage, func() {
		os.RemoveAll(tmpDirectory)
	}
}

func TestStorage_Save(t *testing.T) {
	storage, cleanup := setup(t)
	defer cleanup()

	buf := bytes.NewBuffer([]byte(fileContent))
	writtenBytesNumber, err := storage.Save(savePath, buf)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, writtenBytesNumber, int64(unsafe.Sizeof(buf)))

	file, err := os.Open(filepath.Join(storage.basePath, savePath))
	assert.NoError(t, err)

	data, err := io.ReadAll(file)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, string(data))
}

func TestStorage_Get(t *testing.T) {
	storage, cleanup := setup(t)
	defer cleanup()

	buf := bytes.NewBuffer([]byte(fileContent))
	_, err := storage.Save(savePath, buf)
	assert.NoError(t, err)

	file, err := storage.Get(savePath)
	assert.NoError(t, err)

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	data, err := io.ReadAll(file)
	assert.Equal(t, fileContent, string(data))
}
