package server

import (
	"fmt"

	"github.com/wilmacedo/go-messager/storage"
	"github.com/wilmacedo/go-messager/transport"
)

type Config struct {
	Port         string
	ProducerFunc storage.StorageProducerFunc
}

type Server struct {
	*Config
	topics    map[string]storage.Storage
	consumers []transport.Consumer
	producers []transport.Producer
	exit      chan struct{}
}

func NewServer(cfg *Config) (*Server, error) {
	return &Server{
		Config: cfg,
		topics: make(map[string]storage.Storage),
		exit:   make(chan struct{}),
	}, nil
}

func (s *Server) Start() {
	for _, consumer := range s.consumers {
		if err := consumer.Start(); err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
	}

	for _, producer := range s.producers {
		if err := producer.Start(); err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
	}

	<-s.exit
}
