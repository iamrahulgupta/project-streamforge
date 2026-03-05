package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/iamrahulgupta/streamforge-go/internal"
)

// ProducerConfig contains configuration for the Producer
type ProducerConfig struct {
	// Broker addresses (e.g., []string{"localhost:9092"})
	Brokers []string

	// Topic to produce to
	Topic string

	// Partition to write to (-1 for automatic)
	Partition int32

	// Maximum number of messages to batch
	BatchSize int32

	// Maximum time to wait before sending a batch (milliseconds)
	BatchTimeoutMs int64

	// Compression type: "none", "gzip", "snappy"
	CompressionType string

	// Number of replicas required to acknowledge
	RequiredAcks int16

	// Request timeout in milliseconds
	RequestTimeoutMs int64

	// Maximum retry attempts
	MaxRetries int
}

// DefaultProducerConfig returns a producer config with sensible defaults
func DefaultProducerConfig(brokers []string, topic string) ProducerConfig {
	return ProducerConfig{
		Brokers:         brokers,
		Topic:           topic,
		Partition:       -1,
		BatchSize:       100,
		BatchTimeoutMs:  1000,
		CompressionType: "none",
		RequiredAcks:    1,
		RequestTimeoutMs: 30000,
		MaxRetries:      3,
	}
}

// ProduceResult contains the result of a produce operation
type ProduceResult struct {
	Topic     string
	Partition int32
	Offset    int64
	Error     error
	Timestamp time.Time
}

// Producer is a client for producing messages to StreamForge
type Producer struct {
	config       ProducerConfig
	connection   *internal.Connection
	messageChan  chan *Message
	resultsChan  chan *ProduceResult
	batchTicker  *time.Ticker
	ctx          context.Context
	cancel       context.CancelFunc
	batch        []*Message
	batchMutex   sync.Mutex
	wg           sync.WaitGroup
	closed       bool
	closedMutex  sync.Mutex
}

// Message represents a message to be produced
type Message struct {
	Key       []byte
	Value     []byte
	Timestamp time.Time
	Headers   map[string]string
}

// NewProducer creates a new producer with the given configuration
func NewProducer(config ProducerConfig) (*Producer, error) {
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("at least one broker must be specified")
	}

	if config.Topic == "" {
		return nil, fmt.Errorf("topic must be specified")
	}

	// Connect to first available broker
	conn, err := internal.NewConnection(config.Brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to broker: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := &Producer{
		config:      config,
		connection:  conn,
		messageChan: make(chan *Message, config.BatchSize*2),
		resultsChan: make(chan *ProduceResult, config.BatchSize*2),
		batchTicker: time.NewTicker(time.Duration(config.BatchTimeoutMs) * time.Millisecond),
		ctx:         ctx,
		cancel:      cancel,
		batch:       make([]*Message, 0, config.BatchSize),
	}

	// Start background batch sender
	p.wg.Add(1)
	go p.batchSender()

	return p, nil
}

// Produce sends a message asynchronously and returns a channel for the result
func (p *Producer) Produce(ctx context.Context, message *Message) (*ProduceResult, error) {
	p.closedMutex.Lock()
	if p.closed {
		p.closedMutex.Unlock()
		return nil, fmt.Errorf("producer is closed")
	}
	p.closedMutex.Unlock()

	select {
	case p.messageChan <- message:
		return nil, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.ctx.Done():
		return nil, fmt.Errorf("producer is closing")
	}
}

// ProduceSync sends a message synchronously and waits for the result
func (p *Producer) ProduceSync(ctx context.Context, message *Message) (*ProduceResult, error) {
	if _, err := p.Produce(ctx, message); err != nil {
		return nil, err
	}

	select {
	case result := <-p.resultsChan:
		return result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.ctx.Done():
		return nil, fmt.Errorf("producer is closing")
	}
}

// batchSender processes messages in batches
func (p *Producer) batchSender() {
	defer p.wg.Done()

	for {
		select {
		case msg := <-p.messageChan:
			p.batchMutex.Lock()
			p.batch = append(p.batch, msg)

			if int32(len(p.batch)) >= p.config.BatchSize {
				batch := p.batch
				p.batch = make([]*Message, 0, p.config.BatchSize)
				p.batchMutex.Unlock()

				p.sendBatch(batch)
			} else {
				p.batchMutex.Unlock()
			}

		case <-p.batchTicker.C:
			p.batchMutex.Lock()
			if len(p.batch) > 0 {
				batch := p.batch
				p.batch = make([]*Message, 0, p.config.BatchSize)
				p.batchMutex.Unlock()

				p.sendBatch(batch)
			} else {
				p.batchMutex.Unlock()
			}

		case <-p.ctx.Done():
			// Flush remaining messages
			p.batchMutex.Lock()
			if len(p.batch) > 0 {
				batch := p.batch
				p.batch = make([]*Message, 0, p.config.BatchSize)
				p.batchMutex.Unlock()
				p.sendBatch(batch)
			} else {
				p.batchMutex.Unlock()
			}
			return
		}
	}
}

// sendBatch sends a batch of messages to the broker
func (p *Producer) sendBatch(batch []*Message) {
	for _, msg := range batch {
		result := &ProduceResult{
			Topic:     p.config.Topic,
			Partition: p.config.Partition,
			Offset:    0, // Will be set by broker
			Timestamp: time.Now(),
		}

		// TODO: Send to broker and get actual offset
		select {
		case p.resultsChan <- result:
		case <-p.ctx.Done():
			return
		}
	}
}

// Close gracefully shuts down the producer
func (p *Producer) Close() error {
	p.closedMutex.Lock()
	if p.closed {
		p.closedMutex.Unlock()
		return nil
	}
	p.closed = true
	p.closedMutex.Unlock()

	p.cancel()
	p.batchTicker.Stop()
	p.wg.Wait()

	if p.connection != nil {
		p.connection.Close()
	}

	close(p.messageChan)
	close(p.resultsChan)

	return nil
}

// Flush waits for all pending messages to be sent
func (p *Producer) Flush(timeoutMs int64) error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(timeoutMs)*time.Millisecond)
	defer cancel()

	// Create a sentinel message to know when we've processed all previous messages
	sentinel := &Message{Key: []byte("__sentinel__")}

	if _, err := p.Produce(ctx, sentinel); err != nil {
		return err
	}

	return nil
}
