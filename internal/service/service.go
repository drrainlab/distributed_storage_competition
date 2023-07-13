package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"karma8/internal/storage"
	"log"
)

const storageDir = "storage/"

type Service interface {
	Store(ctx context.Context, name string, size uint64, object io.Reader) error // store file in storage
	Load(ctx context.Context, name string) (io.Reader, error)                    // load file from storage
	Nodes() (capacity uint64, cnt int)                                           // total number of nodes and their capacity
	AddNode(capacity uint64) error                                               // add node to storage
}

type Config struct {
	NodesNum int
	Capacity uint64
}

type filename string
type nodeIdx int

type ObjectStorageService struct {
	cfg     Config
	storage *nodes
}

func NewService(cfg Config) (*ObjectStorageService, error) {
	service := &ObjectStorageService{
		cfg: cfg,
		storage: &nodes{
			list:    make([]storage.Storage, 0),
			binding: make(map[filename][]nodeIdx),
		},
	}
	// calculate single node capacity based on total capacity
	nodeCap := uint64(float64(cfg.Capacity) / float64(cfg.NodesNum))
	// initializing storage nodes
	for i := 0; i < cfg.NodesNum; i++ {
		newNode, err := storage.NewDiskStorage(fmt.Sprintf("node-%d", i), nodeCap, storageDir)
		if err != nil {
			return nil, err
		}
		service.storage.addNode(newNode)
	}

	return service, nil
}

func (s *ObjectStorageService) Nodes() (capacity uint64, cnt int) {
	return s.storage.nodes()
}

func (s *ObjectStorageService) AddNode(capacity uint64) error {
	_, n := s.storage.nodes()
	newNode, err := storage.NewDiskStorage(fmt.Sprintf("node-%d", n+1), capacity, storageDir)
	if err != nil {
		return err
	}
	s.storage.addNode(newNode)
	return nil
}

func (s *ObjectStorageService) Store(ctx context.Context, name string, size uint64, object io.Reader) error {

	if size <= 0 {
		return ErrEmptyFile
	}

	parts, err := s.prepareParts(size)
	if err != nil {
		// // release allocations
		// if len(parts) > 0 {
		// 	for i, part := range parts {
		// 		s.nodes[i].Release(part.ID.String())
		// 	}
		// }
		return err
	}

	reader := bufio.NewReader(object)

	for i := 0; i < len(parts); i++ {
		// dispatch reader that sends data to belonging part data channel
		go func(i int) {
			rcv := parts[i].Receiver()
			for j := 0; j < int(parts[i].Size()); j++ {
				b, err := reader.ReadByte()
				if err != nil && !errors.Is(err, io.EOF) {
					log.Printf("error reading object: %v\n", err)
					return
				}
				rcv <- b
			}
		}(i)
		// start writing data to current node
		if err := s.storage.store(ctx, i, name, parts[i]); err != nil {
			return err
		}
		// waiting for nil on success or error
		err := <-parts[i].Watcher()
		if err != nil {
			return err
		}
		// store ordered list of nodes in map
		s.storage.addBinding(name, i)
	}

	return nil
}

// divide size of file to equally parts, returns file parts objects slice, each index of part belongs to
// corresponding index of node
// (for sake of simplicity we assume that every node has enough capacity for each part)
func (s *ObjectStorageService) prepareParts(size uint64) ([]*storage.FilePart, error) {
	// get nodes number firstly, it can be expanded later
	cap, n := s.Nodes()
	if size > cap {
		return nil, ErrNotEnoughSpace
	}

	var (
		parts           []*storage.FilePart // write will be sequential from beginning, to each node
		mostCapacity    uint64
		mostCapacityIdx int
	)

	partSize := uint64(float64(size) / float64(n))
	// if size too small, just use the first node
	if partSize == 0 {
		n = 1
	}
	residue := size - partSize*uint64(n)

	for i := 0; i < n; i++ {
		// get node capacity
		c := s.storage.getNodeCap(i)

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
		// TODO: release disk space on ctx cancellation or other errors (based on event bus)
		if err := s.storage.alloc(i, part.ID.String(), part.Size()); err != nil {
			return parts, err
		}
		parts = append(parts, part)
	}
	// store residue in store with most space if not equally parted
	if residue > 0 {
		parts[mostCapacityIdx].AdjustSize(residue)
		if err := s.storage.alloc(mostCapacityIdx, parts[mostCapacityIdx].ID.String(), residue); err != nil {
			return parts, err
		}
	}

	return parts, nil
}

func (s *ObjectStorageService) Load(ctx context.Context, name string) (io.Reader, error) {
	nodes := s.storage.getNodes(name)
	var readers []io.Reader

	for _, n := range nodes {
		r, err := n.Get(ctx, name)
		if err != nil {
			return nil, err
		}
		readers = append(readers, r)
	}

	fileReader := io.MultiReader(readers...)

	return fileReader, nil

}
