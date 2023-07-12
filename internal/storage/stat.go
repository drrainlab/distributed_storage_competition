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

func (s *StorageStat) meanCapacity() float64 {
	s.Lock()
	defer s.Unlock()
	return 1 - float64(s.used)/float64(s.total)
}

func (s *StorageStat) decrease(size uint64) {
	s.Lock()
	s.used -= size
	s.Unlock()
}

func (s *StorageStat) increase(size uint64) {
	s.Lock()
	s.used += size
	s.Unlock()
}
