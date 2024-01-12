package server

import (
	"fmt"
	"log/slog"

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

	message chan transport.Message
	exit    chan bool
}

func NewServer(cfg *Config) (*Server, error) {
	message := make(chan transport.Message)

	return &Server{
		Config: cfg,
		topics: make(map[string]storage.Storage),
		producers: []transport.Producer{
			transport.NewHTTPProducer(cfg.Port, message),
		},
		message: message,
		exit:    make(chan bool),
	}, nil
}

func (s *Server) Start() {
	for _, producer := range s.producers {
		go func(p transport.Producer) {
			if err := p.Start(); err != nil {
				slog.Error("%v", err)
			}
		}(producer)
	}

	s.loop()
}

func (s *Server) loop() {
	for {
		select {
		case message := <-s.message:
			if err := s.publish(message); err != nil {
				slog.Error("failed to publish topic", err)
			} else {
				slog.Info("new message produced")
			}
		case <-s.exit:
			fmt.Println("quitting server...")
			return
		}
	}
}

func (s *Server) publish(message transport.Message) error {
	if _, ok := s.topics[message.Topic]; !ok {
		s.topics[message.Topic] = s.ProducerFunc()
		slog.Info("created new topic", "name", message.Topic)
	}

	store := s.topics[message.Topic]

	return store.Push(message.Data)
}
