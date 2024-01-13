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
go run cmd/producer/up.go --topic=new_topic --action=publish
```

## Strategies
* Pub-Sub pattern
* InMemory Storage
* Multi protocols

## Features
* RWMutex to ensure the safety concurrency
* Multi-Thread for each producer/consumer with concurrency
* Fast consumer message serialization
* Up to $62^{62}$ consumers at same time

## Roadmap
- [ ] Authentication to publish/commit and listen
- [ ] Distributed Storage
- [ ] Multiple instance of same application (nodes)