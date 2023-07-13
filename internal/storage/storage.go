package storage

import (
	"context"
	"errors"
	"io"
)

type Storage interface {
	ID() string
	Store(ctx context.Context, name string, part *FilePart) error
	Get(ctx context.Context, name string) (f io.Reader, err error)
	Delete(ctx context.Context, name string) error
	Alloc(id string, size uint64) error // allocate disk space and store allocation by id, if id exists then just increase size
	Capacity() uint64                   // total capacity in bytes
}

var (
	ErrEmptyFile      = errors.New("empty file")
	ErrNotEnoughSpace = errors.New("not enough space")
	ErrAlreadyExist   = errors.New("file already exists")
	ErrNotFound       = errors.New("not found")
)
