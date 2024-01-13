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
		t.Fatal(err)
	}

	fetchedData, err := s.Fetch(0)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(fetchedData, data) {
		t.Fatal("fetched data not match with pushed data")
	}
}

func TestPushAndPop(t *testing.T) {
	s := NewMemoryStorage()
	data := []byte("data")

	err := s.Push(data)
	if err != nil {
		t.Fatal(err)
	}

	err = s.Pop()
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Fetch(1)
	if err == nil {
		t.Fatal("fetched data is need be nil")
	}
}

func TestFetchEmpty(t *testing.T) {
	s := NewMemoryStorage()

	data, err := s.Fetch(0)
	if err == nil {
		t.Fatal("should produce error")
	}

	if data != nil {
		t.Fatal("should produce nil data")
	}
}

func TestOutrangeOffset(t *testing.T) {
	s := NewMemoryStorage()

	err := s.Push([]byte("data"))
	if err != nil {
		t.Fatal("should not produce error")
	}

	data, err := s.Fetch(1)
	if err == nil {
		t.Fatal("should produce error")
	}

	if data != nil {
		t.Fatal("should produce nil data")
	}
}
