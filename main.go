package main

import (
	"log"

	"github.com/wilmacedo/go-messager/server"
	"github.com/wilmacedo/go-messager/storage"
)

func main() {
	cfg := server.Config{
		HTTPPort: ":3000",
		WSPort:   ":4000",
		ProducerFunc: func() storage.Storage {
			return storage.NewMemoryStorage()
		},
	}

	s, err := server.NewServer(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	s.Start()
}
