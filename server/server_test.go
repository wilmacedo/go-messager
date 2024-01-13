package server

import (
	"testing"

	"github.com/wilmacedo/go-messager/storage"
)

func TestNewServer(t *testing.T) {
	cfg := &Config{
		HTTPPort: ":3000",
		WSPort:   ":4000",
		ProducerFunc: func() storage.Storage {
			return storage.NewMemoryStorage()
		},
	}

	s, err := NewServer(cfg)
	if err != nil {
		t.Error(err)
	}

	if s == nil {
		t.Error("server instance should not be nil")
	}
}
