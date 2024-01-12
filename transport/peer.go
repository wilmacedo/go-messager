package transport

import (
	"log/slog"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

type Peer interface {
	Send(b []byte) error
}

type WSPeer struct {
	conn   *websocket.Conn
	server TransportServer
}

func NewWSPeer(conn *websocket.Conn, server TransportServer) *WSPeer {
	p := &WSPeer{
		conn:   conn,
		server: server,
	}

	go p.loop()

	return p
}

func (p *WSPeer) handleHook(h Hook) error {
	if h.Action == "subscribe" {
		p.server.AddPeerToTopic(p, h.Topics)
	}

	return nil
}

func (p *WSPeer) loop() {
	for {
		// TODO: Add retry strategy

		_, incoming, err := p.conn.ReadMessage()
		if err != nil {
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
