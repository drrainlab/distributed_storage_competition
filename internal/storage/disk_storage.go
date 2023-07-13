package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type DiskStorage struct {
	id    string
	path  string
	stat  *StorageStat
	alloc allocks
}

type allocks struct {
	allocations map[string]uint64
	sync.RWMutex
}

func (a *allocks) lock(id string, size uint64) {
	a.Lock()
	defer a.Unlock()
	_, ok := a.allocations[id]
	if ok {
		a.allocations[id] += size
	} else {
		a.allocations[id] = size
	}
}

func (a *allocks) unlock(id string) {
	a.Lock()
	defer a.Unlock()
	if _, ok := a.allocations[id]; !ok {
		return
	}
	delete(a.allocations, id)
}

func (a *allocks) getSize(id string) uint64 {
	a.Lock()
	defer a.Unlock()
	if size, ok := a.allocations[id]; !ok {
		return 0
	} else {
		return size
	}
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

func (s *DiskStorage) Alloc(id string, size uint64) error {
	if s.Capacity() < size {
		return ErrNotEnoughSpace
	}
	s.stat.decrease(size)
	s.alloc.lock(id, size)
	return nil
}

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
		fmt.Println("creating new file name = ", name, " in path = ", path)
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
				log.Printf("ctx finished: %s\n", ctx.Err())
				part.err <- ctx.Err()
				defer s.cleanFailed(path, part)
				return
			case <-part.completed:
				return
			case b := <-part.data:
				fmt.Println("writing ", string(b), " to ", path)
				_, err := f.Write([]byte{b})
				if err != nil {
					defer s.cleanFailed(path, part)
					part.err <- err
					return
				}
				part.count(1)
				if part.size == part.txBytesCnt {
					part.Finish()
				}
			}
		}
	}()

	return nil
}

func (s *DiskStorage) Get(ctx context.Context, name string) (f io.Reader, err error) {
	path := s.buildPath(name)
	// check if file exists
	if _, err := os.Stat(path); !os.IsExist(err) {
		fmt.Println("looking in path: ", path)
		return nil, ErrNotFound
	}

	return os.Open(path)

}

func (s *DiskStorage) Delete(ctx context.Context, name string) error {
	panic("not implemented") // TODO: Implement
}

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
	return fmt.Sprintf("%s/%s", filepath.Dir(s.path)+"/"+s.id, name)
}
