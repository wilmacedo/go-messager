package main

import (
	"flag"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"github.com/wilmacedo/go-messager/transport"
)

func main() {
	topic := flag.String("topic", "test", "topic to execute action")
	action := flag.String("action", "subscribe", "peer action")
	flag.Parse()

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:4000", nil)
	if err != nil {
		panic(err)
	}

	h := transport.Hook{
		Action: *action,
		Topics: []string{*topic},
	}

	data, err := sonic.Marshal(h)
	if err != nil {
		panic(data)
	}

	conn.WriteMessage(websocket.BinaryMessage, data)
	// conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	ch := make(chan any)
	<-ch
}
