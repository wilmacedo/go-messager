package transport

import "github.com/gorilla/websocket"

type Peer interface {
	Send(b []byte) error
}

type WSPeer struct {
	conn *websocket.Conn
}

func NewWSPeer(conn *websocket.Conn) *WSPeer {
	return &WSPeer{
		conn: conn,
	}
}

func (p *WSPeer) Send(b []byte) error {
	return p.conn.WriteMessage(websocket.BinaryMessage, b)
}
