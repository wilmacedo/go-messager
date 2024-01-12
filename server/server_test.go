package server

import (
	"reflect"
	"testing"

	"github.com/wilmacedo/go-messager/storage"
	"github.com/wilmacedo/go-messager/transport"
)

func TestNewServer(t *testing.T) {
	cfg := &Config{
		Port: ":3000",
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

func TestPublish(t *testing.T) {
	cfg := &Config{
		Port: ":3000",
		ProducerFunc: func() storage.Storage {
			return storage.NewMemoryStorage()
		},
	}

	s, err := NewServer(cfg)
	if err != nil {
		t.Error(err)
	}

	message := transport.Message{
		Topic: "test",
		Data:  []byte("testing data"),
	}

	err = s.publish(message)
	if err != nil {
		t.Error(err)
	}

	fetchedData, err := s.topics[message.Topic].Fetch(0)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(fetchedData, message.Data) {
		t.Errorf("fetched data is wrong from message data: got %s want %s", string(fetchedData), string(message.Data))
	}
}
