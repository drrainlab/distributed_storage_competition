package service

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
)

func TestObjectStorageService_Store(t *testing.T) {
	type fields struct {
		cfg Config
	}
	type args struct {
		ctx  context.Context
		name string
		str  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "same content",
			fields: fields{
				cfg: Config{
					NodesNum: 6,
					Capacity: 100000,
				},
			},
			args: args{
				ctx:  context.Background(),
				name: "test_same_content",
				str:  "test_same_content",
			},
		},
		{
			name: "small size",
			fields: fields{
				cfg: Config{
					NodesNum: 6,
					Capacity: 100000,
				},
			},
			args: args{
				ctx:  context.Background(),
				name: "small_size",
				str:  "0",
			},
		},
		{
			name: "bigger size",
			fields: fields{
				cfg: Config{
					NodesNum: 6,
					Capacity: 100000,
				},
			},
			args: args{
				ctx:  context.Background(),
				name: "bigger_size",
				str:  `When doneCh is closed above, the function asyncReader will return. The goroutine we created will also return the next time it evaluates the condition in the for loop. But, what if the goroutine is blocking on r.Read()? Then, we essentially have leaked a goroutine. We’re stuck until the reader unblocks.`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer os.RemoveAll("./storage")

			s, err := NewService(tt.fields.cfg)
			if err != nil {
				t.Errorf("ObjectStorageService init error = %v", err)
			}

			if err := s.Store(tt.args.ctx, tt.args.name, uint64(len(tt.args.str)), strings.NewReader(tt.args.str)); err != nil {
				t.Errorf("ObjectStorageService.Store() error = %v", err)
			}
			reader, err := s.Load(tt.args.ctx, tt.args.name)
			if err != nil {
				t.Errorf("ObjectStorageService.Load() error = %v", err)
			}

			raw, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("reader error = %v", err)
			}

			if string(raw) != tt.args.str {
				t.Errorf("ObjectStorageService.Store() mismatch content")
			}

		})
	}
}

func TestObjectStorageService_StoreAfterNodeAdded(t *testing.T) {
	type fields struct {
		cfg Config
	}
	type args struct {
		ctx  context.Context
		name string
		str  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "added node",
			fields: fields{
				cfg: Config{
					NodesNum: 6,
					Capacity: 100000,
				},
			},
			args: args{
				ctx:  context.Background(),
				name: "test_same_content",
				str:  "test_same_content",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			defer os.RemoveAll("./storage")

			s, err := NewService(tt.fields.cfg)
			if err != nil {
				t.Errorf("ObjectStorageService init error = %v", err)
			}

			if err := s.Store(tt.args.ctx, tt.args.name, uint64(len(tt.args.str)), strings.NewReader(tt.args.str)); err != nil {
				t.Errorf("ObjectStorageService.Store() error = %v", err)
			}

			if err := s.AddNode(10000); err != nil {
				t.Errorf("ObjectStorageService.AddNode() error = %v", err)
			}

			reader, err := s.Load(tt.args.ctx, tt.args.name)
			if err != nil {
				t.Errorf("ObjectStorageService.Load() error = %v", err)
			}

			raw, err := io.ReadAll(reader)
			if err != nil {
				t.Errorf("reader error = %v", err)
			}

			if string(raw) != tt.args.str {
				t.Errorf("ObjectStorageService.Store() mismatch content")
			}

			tt.args.str = tt.args.str + "_after"
			if err := s.Store(tt.args.ctx, tt.args.name+"_after", uint64(len(tt.args.str)), strings.NewReader(tt.args.str)); err != nil {
				t.Errorf("ObjectStorageService.Store() error = %v", err)
			}

			reader, err = s.Load(tt.args.ctx, tt.args.name+"_after")
			if err != nil {
				t.Errorf("ObjectStorageService.Load() error = %v", err)
			}

			raw, err = io.ReadAll(reader)
			if err != nil {
				t.Errorf("reader error = %v", err)
			}

			if string(raw) != tt.args.str {
				t.Errorf("ObjectStorageService.Store() mismatch content")
			}

		})
	}
}
