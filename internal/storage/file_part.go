package storage

import "github.com/google/uuid"

type FilePart struct {
	ID            uuid.UUID
	size          uint64 // total size of the part
	data          chan []byte
	err           chan error // error channel for listening on sending side
	txBytesCnt    uint64     // transferred size in bytes
	isTransferred bool       // flag indicating whether the transfer is completed
}

func NewFilePart(size uint64) *FilePart {
	return &FilePart{
		ID:   uuid.New(),
		size: size,
		data: make(chan []byte),
		err:  make(chan error),
	}
}

func (f *FilePart) count(n uint64) {
	f.txBytesCnt += n
}

func (f *FilePart) Receiver() chan<- []byte {
	return f.data
}

func (f *FilePart) Watcher() <-chan error {
	return f.err
}

func (f *FilePart) Size() uint64 {
	return f.size
}

func (f *FilePart) IsComplete() bool {
	return f.isTransferred
}

func (f *FilePart) Finish() {
	f.isTransferred = true
}
