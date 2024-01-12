package transport

type Message struct {
	Topic string
	Data  []byte
}

type Consumer interface {
	Start() error
}

type Producer interface {
	Start() error
}
