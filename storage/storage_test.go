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
