package backends

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

/*
Local is an implementation of the Storage interface that works with the local disk.
*/
type Local struct {
	maxFileSize uint64 // max number of bytes for files
	basePath    string
}

/*
NewLocal creates a new Local filesystem with the given base path.
basePath: is the base directory to save files to.
maxSize: is the max number of bytes that a file can be.
*/
func NewLocal(basePath string, maxSize uint64) (*Local, error) {
	path, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{basePath: path, maxFileSize: maxSize}, nil
}

/*
fullPath returns the absolute path, appends the given path to the base path.
*/
func (local *Local) fullPath(path string) string {
	return filepath.Join(local.basePath, path)
}

func (local *Local) Save(path string, content io.Reader) (int64, error) {
	fullPath := local.fullPath(path)

	// get the directory and make sure it exists
	uploadPath := filepath.Dir(fullPath)
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return 0, fmt.Errorf("unable to create directory: %w", err)
	}

	// if the file exists then delete it
	_, err = os.Stat(fullPath)
	if err == nil {
		err = os.Remove(fullPath)
		if err != nil {
			return 0, fmt.Errorf("unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return 0, fmt.Errorf("unable to get file info: %w", err)
	}

	bytes := unsafe.Sizeof(content)
	if uint64(bytes) > local.maxFileSize {
		return 0, fmt.Errorf("content size greater than max bytes allowed: %w", err)
	}

	newFile, err := os.Create(fullPath)
	if err != nil {
		return 0, fmt.Errorf("unable to create file: %w", err)
	}
	defer newFile.Close()

	writtenBytes, err := io.Copy(newFile, content)
	if err != nil {
		return 0, fmt.Errorf("unable to write into file: %w", err)
	}

	return writtenBytes, nil
}

func (local *Local) Get(path string) (*os.File, error) {
	fullPath := local.fullPath(path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}

	return file, err
}
