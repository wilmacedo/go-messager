package transport

import (
	"io"
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

	if r.Method == "POST" {
		params := strings.Split(path, "/")
		if len(params) != 2 {
			slog.Error("invalid publish route")
			return
		}

		if params[0] != "publish" {
			slog.Error("invalid action")
			return
		}

		data, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			slog.Error("error reading request body: %v", err)
			return
		}

		topic := params[1]
		p.message <- Message{
			Topic: topic,
			Data:  data,
		}
	}
}

func (p *HTTPProducer) Start() error {
	slog.Info("HTTP transport started", "port", p.Port)
	return http.ListenAndServe(p.Port, p)
}
