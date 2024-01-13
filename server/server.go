package server

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

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
				slog.Error("failed to publish message to topic", err)
				continue
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

	if err := store.Push(message.Data); err != nil {
		return err
	}

	var wg sync.WaitGroup
	var peers []transport.Peer

	start := time.Now()
	counts := make(chan [2]int, len(s.conns))

	for _, peer := range s.conns {
		for _, topic := range peer.GetTopics() {
			if topic == message.Topic {
				peers = append(peers, peer)
			}
		}
	}

	for _, peer := range peers {
		wg.Add(1)

		go func(msg transport.Message, p transport.Peer, c chan [2]int) {
			// TODO: Add retry strategy
			if err := p.Send(msg.Data); err != nil {
				// TODO: Add analytics tracer
				c <- [2]int{0, 1}
				slog.Error("failed to send message to peer", "peer", p.GetID())
			}

			c <- [2]int{1, 0}
			wg.Done()
		}(message, peer, counts)
	}

	go func() {
		wg.Wait()
		close(counts)
	}()

	var total [2]int
	for count := range counts {
		total[0] += count[0]
		total[1] += count[1]
	}

	store.Pop()
	slog.Info("message sended to peers", "took", time.Since(start))
	slog.Info("message status", "success", total[0], "failed", total[1])

	return nil
}

func (s *Server) GetTopics() map[string]storage.Storage {
	return s.topics
}

func (s *Server) RemovePeer(p transport.Peer) error {
	if _, ok := s.conns[p.GetID()]; !ok {
		return fmt.Errorf("peer with id (%s) not found in connection list", p.GetID())
	}

	delete(s.conns, p.GetID())

	return nil
}

func (s *Server) AddPeerToTopic(p transport.Peer, topics []string) error {
	p.AddTopics(topics)
	for _, topic := range topics {
		s.GetStoreByTopic(topic)
	}

	slog.Info("adding new peer to", "topics", topics)
	return nil
}
