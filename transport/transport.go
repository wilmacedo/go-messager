package transport

type Consumer interface {
	Start() error
}

type Producer interface {
	Start() error
}
