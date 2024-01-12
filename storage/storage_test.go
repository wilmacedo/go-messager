package storage

import (
	"reflect"
	"testing"
)

func TestPushAndFetch(t *testing.T) {
	s := NewMemoryStorage()
	data := []byte("data")

	err := s.Push(data)
	if err != nil {
		t.Error(err)
	}

	fetchedData, err := s.Fetch(0)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(fetchedData, data) {
		t.Error("fetched data not match with pushed data")
	}
}

func TestFetchEmpty(t *testing.T) {
	s := NewMemoryStorage()

	data, err := s.Fetch(0)
	if err == nil {
		t.Error("should produce error")
	}

	if data != nil {
		t.Error("should produce nil data")
	}
}

func TestOutrangeOffset(t *testing.T) {
	s := NewMemoryStorage()

	err := s.Push([]byte("data"))
	if err != nil {
		t.Error("should not produce error")
	}

	data, err := s.Fetch(1)
	if err == nil {
		t.Error("should produce error")
	}

	if data != nil {
		t.Error("should produce nil data")
	}
}
