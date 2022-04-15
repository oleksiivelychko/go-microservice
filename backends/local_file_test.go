package backends

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

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
		os.RemoveAll(dir)
	}
}

func TestLocal_Save(t *testing.T) {
	savePath := "/1/test.png"
	fileContent := "Hello World"

	local, cleanup := setup(t)
	defer cleanup()

	err := local.Save(savePath, []byte(fileContent))
	assert.NoError(t, err)

	file, err := os.Open(filepath.Join(local.basePath, savePath))
	assert.NoError(t, err)

	data, err := ioutil.ReadAll(file)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, string(data))
}
