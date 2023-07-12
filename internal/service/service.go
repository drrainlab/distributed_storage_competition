package service

import (
	"context"
	"fmt"
	"io"
	"karma8/internal/storage"
	"strconv"
	"sync"
)

type Service interface {
	Store(ctx context.Context, name string, size uint64, object io.Reader) error
	Load(ctx context.Context, name string) (io.ReadCloser, uint64, error)
	Nodes() (capacity uint64, cnt int) // total number of nodes and their capacity
	AddNode() error
}

type Config struct {
	NodesNum int
	Capacity uint64
}

type ObjectStorageService struct {
	cfg   Config
	nodes []storage.Storage
	sync.RWMutex
}

func NewService(cfg Config) (*ObjectStorageService, error) {
	service := &ObjectStorageService{
		cfg:   cfg,
		nodes: make([]storage.Storage, 0),
	}
	// calculate single node capacity based on total capacity
	nodeCap := uint64(float64(cfg.Capacity) / float64(cfg.NodesNum))
	// initializing storage nodes
	for i := 0; i < cfg.NodesNum; i++ {
		newNode, err := storage.NewDiskStorage(fmt.Sprintf("node-%d", i), nodeCap, "storage/")
		if err != nil {
			return nil, err
		}
		service.nodes = append(service.nodes, newNode)
	}

	return service, nil
}

func (s *ObjectStorageService) Nodes() (capacity uint64, cnt int) {
	s.Lock()
	defer s.Unlock()
	for i, n := range s.nodes {
		cnt = i + 1
		capacity += n.Capacity()
	}
	return
}

func (s *ObjectStorageService) Store(ctx context.Context, name string, size uint64, object io.Reader) error {

	// get nodes number firstly, it can be expanded later
	cap, n := s.Nodes()

	if size > cap {
		return ErrNotEnoughSpace
	}

	var parts []*storage.FilePart // write will be sequential from beginning, to each node

	partSize := uint64(float64(size) / float64(n))

	for i := 0; i < n; i++ {
		c := s.nodes[i].Capacity()
		var part *storage.FilePart
		if c >= partSize {
			part = storage.NewFilePart(partSize)
		} else {
			part = storage.NewFilePart(c)
		}
		// try to allocate disk space on current node
		// on writing failure this allocation will be released
		if err := s.nodes[i].Alloc(part.ID.String(), part.Size()); err != nil {
			return err
		}
		parts = append(parts, part)
	}

	for i := 0; i < len(parts); i++ {

		go func(i int) {
			rcv := parts[i].Receiver()
			rcv <- []byte(strconv.Itoa(i))
			parts[i].Finish()
		}(i)

		if err := s.nodes[i].Store(ctx, name, parts[i]); err != nil {
			return err
		}

		result := <-parts[i].Watcher()
		fmt.Println("result: ", result)

	}

	return nil
}
