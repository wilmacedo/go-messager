# Go Messager
The message queue builded only with standard libraries (and some others :D)

## Usage
To run the application, simple run the following command:
``` make ```

## Strategies
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