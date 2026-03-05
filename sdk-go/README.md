# StreamForge Go SDK

A Go client library for StreamForge, a distributed event streaming platform.

## Features

- **High Performance**: Optimized for throughput and latency
- **Batch Processing**: Automatic message batching for producers
- **Consumer Groups**: Built-in consumer group support with offset management
- **Connection Pooling**: Efficient reuse of broker connections
- **Error Handling**: Robust error handling and recovery
- **Async/Sync API**: Both asynchronous and synchronous APIs

## Installation

```bash
go get github.com/iamrahulgupta/streamforge-go
```

## Quick Start

### Producer

```go
package main

import (
	"context"
	"github.com/iamrahulgupta/streamforge-go/producer"
)

func main() {
	config := producer.DefaultProducerConfig(
		[]string{"localhost:9092"},
		"my-topic",
	)
	
	p, _ := producer.NewProducer(config)
	defer p.Close()
	
	msg := &producer.Message{
		Key:   []byte("key1"),
		Value: []byte("hello world"),
	}
	
	// Async produce
	p.Produce(context.Background(), msg)
	
	// Or sync with confirmation
	result, _ := p.ProduceSync(context.Background(), msg)
	println(result.Offset)
}
```

### Consumer

```go
package main

import (
	"context"
	"github.com/iamrahulgupta/streamforge-go/consumer"
)

func main() {
	config := consumer.DefaultConsumerConfig(
		[]string{"localhost:9092"},
		"my-group",
		[]string{"my-topic"},
	)
	
	c, _ := consumer.NewConsumer(config)
	defer c.Close()
	
	// Poll messages
	msg, _ := c.Poll(context.Background())
	println(string(msg.Value))
	
	// Commit offset
	c.CommitMessage(msg)
}
```

## Project Structure

- **producer/**: Producer implementation
  - `producer.go`: Main producer with batching
  
- **consumer/**: Consumer implementation
  - `consumer.go`: Consumer with group support
  
- **internal/**: Internal utilities
  - `connection.go`: TCP connection management
  - `protocol.go`: Message serialization/deserialization
  
- **examples/**: Example applications
  - `producer_example.go`: Producer usage examples
  - `consumer_example.go`: Consumer usage examples

## Configuration

### Producer Configuration

```go
config := producer.ProducerConfig{
    Brokers:         []string{"localhost:9092"},
    Topic:           "events",
    Partition:       -1,  // Auto partition
    BatchSize:       100,
    BatchTimeoutMs:  1000,
    CompressionType: "none",
    RequiredAcks:    1,
    RequestTimeoutMs: 30000,
    MaxRetries:      3,
}
```

### Consumer Configuration

```go
config := consumer.ConsumerConfig{
    Brokers:             []string{"localhost:9092"},
    GroupID:             "my-group",
    Topics:              []string{"topic1", "topic2"},
    InitialOffset:       "latest",  // or "earliest"
    SessionTimeoutMs:    10000,
    HeartbeatIntervalMs: 3000,
    MaxPollRecords:      500,
    FetchMinBytes:       1,
    FetchMaxBytes:       52428800,
    FetchTimeoutMs:      30000,
}
```

## API Documentation

### Producer

#### NewProducer(config ProducerConfig) (*Producer, error)
Creates a new producer instance.

#### Produce(ctx context.Context, message *Message) (*ProduceResult, error)
Sends a message asynchronously.

#### ProduceSync(ctx context.Context, message *Message) (*ProduceResult, error)
Sends a message synchronously and waits for confirmation.

#### Flush(timeoutMs int64) error
Waits for all pending messages to be sent.

#### Close() error
Gracefully shuts down the producer.

### Consumer

#### NewConsumer(config ConsumerConfig) (*Consumer, error)
Creates a new consumer instance.

#### Poll(ctx context.Context) (*ConsumerMessage, error)
Retrieves a single message.

#### PollMessages(ctx context.Context, maxMessages int) ([]*ConsumerMessage, error)
Retrieves multiple messages in one call.

#### CommitOffset(topic string, partition int32, offset int64) error
Commits an offset for a topic-partition.

#### CommitMessage(msg *ConsumerMessage) error
Commits the offset of a consumed message.

#### Subscribe(topics []string) error
Subscribes to additional topics.

#### Close() error
Gracefully shuts down the consumer.

## Examples

See the `examples/` directory for complete working examples:

### Run Producer Example
```bash
cd examples
go run producer_example.go
```

### Run Consumer Example
```bash
cd examples
go run consumer_example.go
```

## Testing

```bash
go test ./...
go test -v ./producer
go test -v ./consumer
```

## Benchmarks

```bash
go test -bench=. -benchmem ./producer
go test -bench=. -benchmem ./consumer
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For issues, questions, or suggestions:
- Open an issue on GitHub
- Check existing issues and discussions
- Review the examples in the `examples/` directory

## Roadmap

- [ ] Schema Registry integration
- [ ] Exactly-once semantics
- [ ] Transactions support
- [ ] Metrics and monitoring
- [ ] Compression codecs (gzip, snappy, lz4)
- [ ] Performance optimizations
- [ ] Additional examples and tutorials
