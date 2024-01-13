# Go Messager
The message queue builded only with standard libraries (and some others :D)

## Usage
To run the application server, simple run the following command:
```bash
make 
```

To run single consumer, run the following command:
```bash
make up-consumer 
```

To change the topic or message, run the following command:
```bash
go run cmd/producer/up.go --topic=new_topic --action=subscribe
```

To publish data in a topic, run the following command:
```bash
curl -X POST http://localhost:3000/publish/topic --data-binary 'any data'
```

## Strategies
* Pub-Sub pattern
* InMemory Storage
* Multi protocols

## Features
* RWMutex to ensure the safety concurrency
* Concurrency to speed data communication for each consumer
* Fast message serialization with [sonic](https://github.com/bytedance/sonic)
* Up to $62^{62}$ consumers at same time

## Roadmap
- [ ] Benchmark for performance at high concurrency
- [ ] Authentication to publish/commit and listen
- [ ] Distributed Storage
- [ ] Multiple instance of same application (nodes)