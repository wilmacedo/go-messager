run: build
	@./bin/go-messager

build:
	@go build -o bin/go-messager

test:
	@go test ./...

up-consumer:
	@go run cmd/consumer/up.go