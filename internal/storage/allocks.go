package storage

import "sync"

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
