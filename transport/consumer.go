package transport

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader

type Hook struct {
	Action string   `json:"action"`
	Topics []string `json:"topics"`
}

type WSConsumer struct {
	Port   string
	Peer   chan<- Peer
	Server TransportServer
}

func NewWSConsumer(port string, peer chan<- Peer, server TransportServer) *WSConsumer {
	return &WSConsumer{
		Port:   port,
		Peer:   peer,
		Server: server,
	}
}

func (ws *WSConsumer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	p, err := NewWSPeer(conn, ws.Server)
	if err != nil {
		fmt.Println(err)
		return
	}

	ws.Peer <- p
}

func (ws *WSConsumer) Start() error {
	slog.Info("WS transport started", "port", ws.Port)
	return http.ListenAndServe(ws.Port, ws)
}
