package reader

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ReadFile(path string) ([]byte, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return []byte{}, err
	}

	if _, err = os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return []byte{}, fmt.Errorf("path %s is invalid", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	return io.ReadAll(file)
}
