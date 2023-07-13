package service

import (
	"context"
	"karma8/internal/storage"
	"sync"
)

type nodes struct {
	list    []storage.Storage
	binding map[filename][]nodeIdx // storing here files <-> nodes bindings
	mu      sync.RWMutex
}

func (s *nodes) nodes() (capacity uint64, cnt int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, n := range s.list {
		cnt = i + 1
		capacity += n.Capacity()
	}
	return
}

func (s *nodes) getNodeCap(i int) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.list[i].Capacity()
}

func (s *nodes) addNode(n storage.Storage) {
	s.mu.Lock()
	s.list = append(s.list, n)
	s.mu.Unlock()
}

func (s *nodes) addBinding(name string, idx int) {
	s.mu.Lock()
	s.binding[filename(name)] = append(s.binding[filename(name)], nodeIdx(idx))
	s.mu.Unlock()
}

func (s *nodes) getNodes(name string) (bindedNodes []storage.Storage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, n := range s.binding[filename(name)] {
		bindedNodes = append(bindedNodes, s.list[n])
	}
	return
}

func (s *nodes) alloc(nodeIdx int, id string, size uint64) error {
	return s.list[nodeIdx].Alloc(id, size)
}

func (s *nodes) store(ctx context.Context, nodeIdx int, name string, part *storage.FilePart) error {
	return s.list[nodeIdx].Store(ctx, name, part)
}
