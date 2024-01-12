package transport

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type WSConsumer struct {
	Port string
	Peer chan<- Peer
}

func NewWSConsumer(port string, peer chan<- Peer) *WSConsumer {
	return &WSConsumer{
		Port: port,
		Peer: peer,
	}
}

func (ws *WSConsumer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	p := NewWSPeer(conn)
	ws.Peer <- p
}

func (ws *WSConsumer) Start() error {
	slog.Info("WS transport started", "port", ws.Port)
	return http.ListenAndServe(ws.Port, ws)
}
