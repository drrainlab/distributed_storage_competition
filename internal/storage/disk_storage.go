package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type DiskStorage struct {
	id    string
	path  string
	stat  *StorageStat
	alloc allocks
}

func NewDiskStorage(id string, capacity uint64, path string) (Storage, error) {
	if err := os.MkdirAll(filepath.Dir(path)+"/"+id, os.ModePerm); err != nil {
		return nil, err
	}

	return &DiskStorage{
		id: id, path: path, stat: &StorageStat{
			used:  0,
			total: capacity,
		},
		alloc: allocks{
			allocations: make(map[string]uint64),
		},
	}, nil
}

func (s *DiskStorage) ID() string {
	return s.id
}

// allocates space on disk (current node)
func (s *DiskStorage) Alloc(id string, size uint64) error {
	if s.Capacity() < size {
		return ErrNotEnoughSpace
	}
	s.stat.decrease(size)
	s.alloc.lock(id, size)
	return nil
}

// channel based file uploading
func (s *DiskStorage) Store(ctx context.Context, name string, part *FilePart) error {

	switch {
	case part.size == 0:
		return ErrEmptyFile
	case s.stat.capacity() < part.size:
		return ErrNotEnoughSpace
	}
	path := s.buildPath(name)
	// check if file exists
	if _, err := os.Stat(path); os.IsExist(err) {
		return ErrAlreadyExist
	}

	// start writing in background
	go func() {
		f, err := os.Create(path)
		if err != nil {
			part.err <- err
			return
		}
		defer func() {
			close(part.data)
			close(part.err)
			if err := f.Close(); err != nil {
				log.Printf("closing file error: %s\n", err.Error())
			}
		}()

		for {
			select {
			case <-ctx.Done():
				part.err <- ctx.Err()
				defer s.cleanFailed(path, part)
				return
			case <-part.completed:
				return
			case b := <-part.data:
				_, err := f.Write([]byte{b})
				if err != nil {
					defer s.cleanFailed(path, part)
					part.err <- err
					return
				}
				part.count(1)
				// if bytes transferred is equal to the size, then finish writing
				if part.size == part.txBytesCnt {
					part.Finish()
				}
			}
		}
	}()

	return nil
}

// returns file io.reader or error
func (s *DiskStorage) Get(ctx context.Context, name string) (f io.Reader, err error) {
	path := s.buildPath(name)
	// check if file exists
	if _, err := os.Stat(path); err != nil && !os.IsExist(err) {
		fmt.Println("looking in path: ", path, err)
		return nil, ErrNotFound
	}

	// NOTE: didn't make cancelable reader for simplicity
	return os.Open(path)
}

func (s *DiskStorage) Delete(ctx context.Context, name string) error {
	panic("not implemented") // TODO: Implement
}

// capacity of node
func (s *DiskStorage) Capacity() uint64 {
	return s.stat.capacity()
}

// deletes old file and
func (s *DiskStorage) cleanFailed(path string, part *FilePart) {
	id := part.ID.String()
	size := s.alloc.getSize(id)
	// release allocated size
	s.stat.increase(size)
	// remove allocation
	s.alloc.unlock(part.ID.String())
	if err := os.Remove(path); err != nil {
		log.Printf("file removing error: %s\n", err.Error())
	}
}

func (s *DiskStorage) buildPath(name string) string {
	return fmt.Sprintf("./%s/%s", filepath.Dir(s.path)+"/"+s.id, name)
}
