package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

/*
Local is an implementation of the storage.ILocal interface that works with the local disk.
*/
type Local struct {
	maxFileSize uint64 // max number of bytes for files
	basePath    string
}

/*
New creates a new Local file system with the given base path.
basePath: is the base directory to save files to.
maxSize: is the max number of bytes that a file can be.
*/
func New(basePath string, maxSize uint64) (*Local, error) {
	path, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{basePath: path, maxFileSize: maxSize}, nil
}

/*
getBasePath returns the absolute path, appends the given path to the base path.
*/
func (storage *Local) getBasePath(path string) string {
	return filepath.Join(storage.basePath, path)
}

func (storage *Local) Save(path string, data io.Reader) (int64, error) {
	basePath := storage.getBasePath(path)

	// make sure the directory already exists
	uploadPath := filepath.Dir(basePath)
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return 0, err
	}

	// if the file exists, delete it
	_, err = os.Stat(basePath)
	if err == nil {
		err = os.Remove(basePath)
		if err != nil {
			return 0, err
		}
	} else if !os.IsNotExist(err) {
		return 0, fmt.Errorf("unable to fetch info for %s: %s", basePath, err)
	}

	bytes := unsafe.Sizeof(data)
	if uint64(bytes) > storage.maxFileSize {
		return 0, fmt.Errorf("bytes size %d greater than bytes %d is allowed: %s", uint64(bytes), storage.maxFileSize, err)
	}

	file, err := os.Create(basePath)
	if err != nil {
		return 0, err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	writtenBytesNumber, err := io.Copy(file, data)
	if err != nil {
		return 0, err
	}

	return writtenBytesNumber, nil
}

func (storage *Local) Get(path string) (*os.File, error) {
	return os.Open(storage.getBasePath(path))
}
