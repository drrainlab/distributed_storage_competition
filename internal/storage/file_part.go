package storage

import "github.com/google/uuid"

type FilePart struct {
	ID         uuid.UUID
	size       uint64 // total size of the part
	data       chan byte
	err        chan error // error channel for listening on sending side
	txBytesCnt uint64     // transferred size in bytes
	completed  chan bool
}

func NewFilePart(size uint64) *FilePart {
	return &FilePart{
		ID:        uuid.New(),
		size:      size,
		data:      make(chan byte, 100000),
		err:       make(chan error, 1),
		completed: make(chan bool, 1),
	}
}

func (f *FilePart) count(n uint64) {
	f.txBytesCnt += n
}

func (f *FilePart) Receiver() chan<- byte {
	return f.data
}

func (f *FilePart) Watcher() <-chan error {
	return f.err
}

func (f *FilePart) Size() uint64 {
	return f.size
}

func (f *FilePart) AdjustSize(s uint64) {
	f.size += s
}

func (f *FilePart) Finish() {
	f.completed <- true
}
