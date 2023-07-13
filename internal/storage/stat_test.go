package storage

import (
	"sync"
	"testing"
)

func TestStorageStat_capacity(t *testing.T) {
	type fields struct {
		used    uint64
		total   uint64
		RWMutex sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "capacity",
			fields: fields{
				used:  100,
				total: 1000,
			},
			want: 900,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StorageStat{
				used:    tt.fields.used,
				total:   tt.fields.total,
				RWMutex: tt.fields.RWMutex,
			}
			if got := s.capacity(); got != tt.want {
				t.Errorf("StorageStat.capacity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageStat_increase(t *testing.T) {
	type fields struct {
		used    uint64
		total   uint64
		RWMutex sync.RWMutex
	}
	type args struct {
		size uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "increase",
			fields: fields{
				used:  100,
				total: 1000,
			},
			args: args{
				size: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StorageStat{
				used:    tt.fields.used,
				total:   tt.fields.total,
				RWMutex: tt.fields.RWMutex,
			}
			s.increase(tt.args.size)
			cap := s.capacity()
			if cap != 1000 {
				t.Errorf("increase failed: wanted %d have %d", 1000, cap)
			}
		})
	}
}

func TestStorageStat_decrease(t *testing.T) {
	type fields struct {
		used    uint64
		total   uint64
		RWMutex sync.RWMutex
	}
	type args struct {
		size uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "decrease",
			fields: fields{
				used:  100,
				total: 1000,
			},
			args: args{
				size: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StorageStat{
				used:    tt.fields.used,
				total:   tt.fields.total,
				RWMutex: tt.fields.RWMutex,
			}
			s.decrease(tt.args.size)
			cap := s.capacity()
			if cap != 800 {
				t.Errorf("increase failed: wanted %d have %d", 800, cap)
			}
		})
	}
}
