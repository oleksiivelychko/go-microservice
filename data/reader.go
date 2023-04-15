package data

import (
	"fmt"
	"github.com/oleksiivelychko/go-microservice/utils"
	"io"
	"os"
	"path/filepath"
)

func ReadFile(path string) ([]byte, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return []byte{}, err
	}

	if !utils.IsPathValid(path) {
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
