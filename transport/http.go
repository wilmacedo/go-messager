package transport

import (
	"log/slog"
	"net/http"
	"strings"
)

type HTTPProducer struct {
	Port    string
	message chan<- Message
}

func NewHTTPProducer(port string, message chan Message) *HTTPProducer {
	return &HTTPProducer{
		Port:    port,
		message: message,
	}
}

func (p *HTTPProducer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	if r.Method == "GET" {
		// commit
	}

	if r.Method == "POST" {
		params := strings.Split(path, "/")
		if len(params) != 2 {
			slog.Error("invalid publish route")
			return
		}

		topic := params[1]
		p.message <- Message{
			Topic: topic,
			Data:  []byte("hello world"),
		}
	}
}

func (p *HTTPProducer) Start() error {
	slog.Info("HTTP transport started", "port", p.Port)
	return http.ListenAndServe(p.Port, p)
}
