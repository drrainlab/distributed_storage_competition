package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"karma8/internal/storage"
	"log"
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

// divide size of file to equally parts, returns file parts objects slice, each index of part belongs to
// corresponding index of node
func (s *ObjectStorageService) prepareParts(size uint64) ([]*storage.FilePart, error) {
	// get nodes number firstly, it can be expanded later
	cap, n := s.Nodes()

	if size > cap {
		return nil, ErrNotEnoughSpace
	}

	var parts []*storage.FilePart // write will be sequential from beginning, to each node

	partSize := uint64(float64(size) / float64(n))
	residue := size % partSize

	var (
		mostCapacity    uint64
		mostCapacityIdx int
	)

	for i := 0; i < n; i++ {
		// get node capacity
		c := s.nodes[i].Capacity()

		var part *storage.FilePart

		if c >= mostCapacity {
			mostCapacity = c
			mostCapacityIdx = i
		}

		if c >= partSize {
			part = storage.NewFilePart(partSize)
		} else {
			part = storage.NewFilePart(c)
		}
		// try to allocate disk space on current node
		// on writing failure this allocation will be released
		if err := s.nodes[i].Alloc(part.ID.String(), part.Size()); err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	// store residue in store with most space if not equally parted
	if residue > 0 {
		parts[mostCapacityIdx].AdjustSize(residue)
		s.nodes[mostCapacityIdx].Alloc(parts[mostCapacityIdx].ID.String(), residue)
	}

	return parts, nil
}

func (s *ObjectStorageService) Store(ctx context.Context, name string, size uint64, object io.Reader) error {

	parts, err := s.prepareParts(size)
	if err != nil {
		return err
	}

	// capcap, nn := s.Nodes()
	// fmt.Println("nodes after alloc: ", capcap, nn)

	reader := bufio.NewReader(object)

	for i := 0; i < len(parts); i++ {

		// dispatch reader that sends data to belonging part data channel
		go func(i int) {
			rcv := parts[i].Receiver()
			for j := 0; j < int(parts[i].Size()); j++ {
				b, err := reader.ReadByte()
				// fmt.Println("read ", string(b), " from source", " part=", i, " size=", parts[i].Size())
				if err != nil && !errors.Is(err, io.EOF) {
					log.Printf("error reading object: %v\n", err)
					return
				}
				rcv <- b
			}
		}(i)

		// start writing data to current node
		if err := s.nodes[i].Store(ctx, name, parts[i]); err != nil {
			return err
		}

		// waiting for nil on success or error
		err := <-parts[i].Watcher()
		if err != nil {
			return err
		}

	}

	return nil
}
