package backends

import (
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"unsafe"
)

/*
Local is an implementation of the Storage interface which works with the local disk.
*/
type Local struct {
	maxFileSize uint64 // max number of bytes for files
	basePath    string
}

/*
NewLocal creates a new Local filesystem with the given base path
basePath: is the base directory to save files to
maxSize: is the max number of bytes that a file can be
*/
func NewLocal(basePath string, maxSize uint64) (*Local, error) {
	path, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{basePath: path, maxFileSize: maxSize}, nil
}

/*
fullPath returns the absolute path
*/
func (l *Local) fullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(l.basePath, path)
}

func (l *Local) Save(path string, content []byte) error {
	fullPath := l.fullPath(path)

	// get the directory and make sure it exists
	uploadPath := filepath.Dir(fullPath)
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("unable to create directory: %w", err)
	}

	// if the file exists delete it
	_, err = os.Stat(fullPath)
	if err == nil {
		err = os.Remove(fullPath)
		if err != nil {
			return xerrors.Errorf("unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("unable to get file info: %w", err)
	}

	bytes := unsafe.Sizeof(content)
	if uint64(bytes) > l.maxFileSize {
		return xerrors.Errorf("content size greater than max bytes allowed: %w", err)
	}

	err = os.WriteFile(fullPath, content, 0644)
	if err != nil {
		return xerrors.Errorf("unable to write to file: %w", err)
	}

	return nil
}

func (l *Local) Get(path string) (*os.File, error) {
	fullPath := l.fullPath(path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, xerrors.Errorf("unable to open file: %w", err)
	}

	return file, err
}
