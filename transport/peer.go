package transport

import (
	"log/slog"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/wilmacedo/go-messager/util"
)

type Peer interface {
	Send(b []byte) error
	GetID() string
	AddTopics(topics []string)
	GetTopics() []string
}

type WSPeer struct {
	conn   *websocket.Conn
	server TransportServer
	id     string
	topics []string
}

func NewWSPeer(conn *websocket.Conn, server TransportServer) (*WSPeer, error) {
	id, err := util.GenerateRandomID()
	if err != nil {
		return nil, err
	}

	p := &WSPeer{
		conn:   conn,
		server: server,
		id:     id,
		topics: make([]string, 0),
	}

	go p.loop()

	return p, nil
}

func (p *WSPeer) handleHook(h Hook) error {
	if h.Action == "subscribe" {
		p.server.AddPeerToTopic(p, h.Topics)
	}

	return nil
}

func (p *WSPeer) loop() {
	for {
		_, incoming, err := p.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				if err := p.server.RemovePeer(p); err != nil {
					slog.Error("ws peer remove error", "err", err)
					return
				}

				slog.Info("ws peer disconnected")
				return
			}

			slog.Error("ws peer read message error", "err", err)
			return
		}

		var h Hook
		if err := sonic.Unmarshal(incoming, &h); err != nil {
			slog.Error("ws peer unmarshal error", "err", err)
			return
		}

		if err := p.handleHook(h); err != nil {
			slog.Error("ws peer handle error", "err", err)
			return
		}
	}
}

func (p *WSPeer) Send(b []byte) error {
	return p.conn.WriteMessage(websocket.BinaryMessage, b)
}

func (p *WSPeer) GetID() string {
	return p.id
}

func (p *WSPeer) AddTopics(topics []string) {
	for _, topic := range topics {
		has := false
		for _, t := range p.topics {
			if t == topic {
				has = true
			}
		}

		if !has {
			p.topics = append(p.topics, topic)
		}
	}
}

func (p *WSPeer) GetTopics() []string {
	return p.topics
}
