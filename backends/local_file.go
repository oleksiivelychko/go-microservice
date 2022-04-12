package backends

import (
	"golang.org/x/xerrors"
	"io"
	"os"
	"path/filepath"
)

/*
Local is an implementation of the Storage interface which works with the local disk.
*/
type Local struct {
	maxFileSize int // max number of bytes for files
	basePath    string
}

/*
NewLocal creates a new Local filesystem with the given base path
basePath: is the base directory to save files to
maxSize: is the max number of bytes that a file can be
*/
func NewLocal(basePath string, maxSize int) (*Local, error) {
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

func (l *Local) Save(path string, contents io.Reader) error {
	fullPath := l.fullPath(path)

	// get the directory and make sure it exists
	uploadDir := filepath.Dir(fullPath)
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("unable to create directory: %w", err)
	}

	// if the file exists delete it
	_, err = os.Stat(uploadDir)
	if err == nil {
		err = os.Remove(uploadDir)
		if err != nil {
			return xerrors.Errorf("unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("unable to get file info: %w", err)
	}

	f, err := os.Create(uploadDir)
	if err != nil {
		return xerrors.Errorf("unable to create file: %w", err)
	}

	defer f.Close()

	// write the contents to the new file
	// TODO: ensure that we are not writing greater than max bytes
	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("unable to write to file: %w", err)
	}

	return nil
}

func (l *Local) Get(path string) (*os.File, error) {
	fullPath := l.fullPath(path)

	f, err := os.Open(fullPath)
	if err != nil {
		return nil, xerrors.Errorf("unable to open file: %w", err)
	}

	return f, err
}
