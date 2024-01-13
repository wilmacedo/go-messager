package transport

import "github.com/wilmacedo/go-messager/storage"

type Message struct {
	Topic string
	Data  []byte
}

type TransportServer interface {
	GetTopics() map[string]storage.Storage
	AddPeerToTopic(p Peer, topics []string) error
	RemovePeer(p Peer) error
	GetStoreByTopic(topic string) storage.Storage
}

type Consumer interface {
	Start() error
}

type Producer interface {
	Start() error
}
