package storage

import "sync"

type StorageStat struct {
	used  uint64
	total uint64
	sync.RWMutex
}

func (s *StorageStat) capacity() uint64 {
	s.Lock()
	defer s.Unlock()
	return s.total - s.used
}

// decrease storage capacity
func (s *StorageStat) decrease(size uint64) {
	s.Lock()
	s.used += size
	s.Unlock()
}

// increase storage capacity
func (s *StorageStat) increase(size uint64) {
	s.Lock()
	s.used -= size
	s.Unlock()
}
