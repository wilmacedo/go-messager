package server

import (
	"fmt"
	"log/slog"

	"github.com/wilmacedo/go-messager/storage"
	"github.com/wilmacedo/go-messager/transport"
)

type Config struct {
	HTTPPort     string
	WSPort       string
	ProducerFunc storage.StorageProducerFunc
}

type Server struct {
	*Config

	topics    map[string]storage.Storage
	producers []transport.Producer
	conns     map[string]transport.Peer
	consumers []transport.Consumer

	peer    chan transport.Peer
	message chan transport.Message
	exit    chan bool
}

func NewServer(cfg *Config) (*Server, error) {
	message := make(chan transport.Message)
	peer := make(chan transport.Peer)

	s := &Server{
		Config: cfg,
		topics: make(map[string]storage.Storage),
		conns:  make(map[string]transport.Peer),
		producers: []transport.Producer{
			transport.NewHTTPProducer(cfg.HTTPPort, message),
		},
		peer:    peer,
		message: message,
		exit:    make(chan bool),
	}

	s.consumers = []transport.Consumer{
		transport.NewWSConsumer(cfg.WSPort, peer, s),
	}

	return s, nil
}

func (s *Server) Start() {
	for _, consumer := range s.consumers {
		go func(c transport.Consumer) {
			if err := c.Start(); err != nil {
				slog.Error("%v", err)
			}
		}(consumer)
	}

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
		case peer := <-s.peer:
			slog.Info("new peer connected")
			s.conns[peer.GetID()] = peer
		case <-s.exit:
			fmt.Println("quitting server...")
			return
		}
	}
}

func (s *Server) GetStoreByTopic(topic string) storage.Storage {
	if _, ok := s.topics[topic]; !ok {
		s.topics[topic] = s.ProducerFunc()
		slog.Info("created new topic", "name", topic)
	}

	return s.topics[topic]
}

func (s *Server) publish(message transport.Message) error {
	store := s.GetStoreByTopic(message.Topic)

	return store.Push(message.Data)
}

func (s *Server) GetTopics() map[string]storage.Storage {
	return s.topics
}

func (s *Server) RemovePeer(p transport.Peer) error {
	if _, ok := s.conns[p.GetID()]; !ok {
		return fmt.Errorf("peer with id (%s) not found in connection list", p.GetID())
	}

	delete(s.conns, p.GetID())

	fmt.Println(s.conns)

	return nil
}

func (s *Server) AddPeerToTopic(p transport.Peer, topics []string) error {
	for _, topic := range topics {
		store := s.GetStoreByTopic(topic)

		size := store.Size()
		for i := 0; i < size; i++ {
			b, err := store.Fetch(uint(i))
			if err != nil {
				slog.Error("offset not founded in topic store", "offset", i)
				continue
			}

			fmt.Println(string(b))
		}
	}

	slog.Info("adding new peer to", "topics", topics)
	return nil
}
