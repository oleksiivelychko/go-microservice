package storage

import "io"

type ILocal interface {
	Save(path string, file io.Reader) (int64, error)
}
