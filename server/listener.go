package server

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/wilmacedo/go-messager/transport"
)

type Listener struct {
	server *Server
}

func NewListener(server *Server) *Listener {
	return &Listener{
		server: server,
	}
}

func (l *Listener) Listen() {
	for {
		select {
		case message := <-l.server.message:
			if err := l.publish(message); err != nil {
				slog.Error("failed to publish message to topic", err)
				continue
			}
		case peer := <-l.server.peer:
			slog.Info("new peer connected")
			l.server.conns[peer.GetID()] = peer
		case <-l.server.exit:
			fmt.Println("quitting server...")
			return
		}
	}
}

func (l *Listener) publish(message transport.Message) error {
	store := l.server.GetStoreByTopic(message.Topic)

	if err := store.Push(message.Data); err != nil {
		return err
	}

	var wg sync.WaitGroup
	var peers []transport.Peer

	start := time.Now()
	counts := make(chan [2]int, len(l.server.conns))

	for _, peer := range l.server.conns {
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
