package storage

import (
	"fmt"
	"sync"
)

type Storage interface {
	Push([]byte) error
	Pop() error
	Fetch(uint) ([]byte, error)
	Size() int
}

type StorageProducerFunc func() Storage

type MemoryStorage struct {
	mu   sync.RWMutex
	data [][]byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make([][]byte, 0),
	}
}

func (s *MemoryStorage) Push(b []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, b)
	return nil
}

func (s *MemoryStorage) Pop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.data) == 0 {
		return nil // return empty error?
	}

	s.data = s.data[:len(s.data)-1]
	return nil
}

func (s *MemoryStorage) Fetch(offset uint) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	size := len(s.data)

	if size == 0 {
		return nil, fmt.Errorf("storage is empty")
	}

	if size-1 < int(offset) {
		return nil, fmt.Errorf("offset (%d) out of range", offset)
	}

	fetch := s.data[offset]
	if fetch == nil {
		return nil, fmt.Errorf("missing bytes in offset (%d)", offset)
	}

	return fetch, nil
}

func (s *MemoryStorage) Size() int {
	return len(s.data)
}
