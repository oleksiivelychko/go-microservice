package contracts

import "io"

type Storage interface {
	Save(path string, file io.Reader) error
}
