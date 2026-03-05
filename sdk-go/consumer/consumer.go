package consumer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/iamrahulgupta/streamforge-go/internal"
)

// ConsumerConfig contains configuration for the Consumer
type ConsumerConfig struct {
	// Broker addresses (e.g., []string{"localhost:9092"})
	Brokers []string

	// Consumer group ID
	GroupID string

	// Topics to subscribe to
	Topics []string

	// Initial offset: "earliest" or "latest"
	InitialOffset string

	// Session timeout in milliseconds
	SessionTimeoutMs int64

	// Heartbeat interval in milliseconds
	HeartbeatIntervalMs int64

	// Max poll records per fetch
	MaxPollRecords int32

	// Fetch minimum bytes
	FetchMinBytes int32

	// Fetch maximum bytes
	FetchMaxBytes int32

	// Fetch timeout in milliseconds
	FetchTimeoutMs int64
}

// DefaultConsumerConfig returns a consumer config with sensible defaults
func DefaultConsumerConfig(brokers []string, groupID string, topics []string) ConsumerConfig {
	return ConsumerConfig{
		Brokers:             brokers,
		GroupID:             groupID,
		Topics:              topics,
		InitialOffset:       "latest",
		SessionTimeoutMs:    10000,
		HeartbeatIntervalMs: 3000,
		MaxPollRecords:      500,
		FetchMinBytes:       1,
		FetchMaxBytes:       52428800, // 50MB
		FetchTimeoutMs:      30000,
	}
}

// ConsumerMessage represents a message received from the broker
type ConsumerMessage struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
	Timestamp time.Time
	Headers   map[string]string
}

// Consumer is a client for consuming messages from StreamForge
type Consumer struct {
	config          ConsumerConfig
	connection      *internal.Connection
	messagesChan    chan *ConsumerMessage
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	closed          bool
	closedMutex     sync.Mutex
	currentOffsets  map[string]map[int32]int64 // topic -> partition -> offset
	offsetsMutex    sync.RWMutex
	heartbeatTicker *time.Ticker
}

// NewConsumer creates a new consumer with the given configuration
func NewConsumer(config ConsumerConfig) (*Consumer, error) {
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("at least one broker must be specified")
	}

	if config.GroupID == "" {
		return nil, fmt.Errorf("group ID must be specified")
	}

	if len(config.Topics) == 0 {
		return nil, fmt.Errorf("at least one topic must be specified")
	}

	// Connect to first available broker
	conn, err := internal.NewConnection(config.Brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to broker: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &Consumer{
		config:         config,
		connection:     conn,
		messagesChan:   make(chan *ConsumerMessage, config.MaxPollRecords),
		ctx:            ctx,
		cancel:         cancel,
		currentOffsets: make(map[string]map[int32]int64),
		heartbeatTicker: time.NewTicker(
			time.Duration(config.HeartbeatIntervalMs) * time.Millisecond,
		),
	}

	// Initialize offset map
	for _, topic := range config.Topics {
		c.currentOffsets[topic] = make(map[int32]int64)
	}

	// Start background fetcher and heartbeat
	c.wg.Add(2)
	go c.messageFetcher()
	go c.heartbeatSender()

	return c, nil
}

// Poll retrieves messages from the consumer
// Context timeout controls how long to wait for messages
func (c *Consumer) Poll(ctx context.Context) (*ConsumerMessage, error) {
	c.closedMutex.Lock()
	if c.closed {
		c.closedMutex.Unlock()
		return nil, fmt.Errorf("consumer is closed")
	}
	c.closedMutex.Unlock()

	select {
	case msg := <-c.messagesChan:
		return msg, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.ctx.Done():
		return nil, fmt.Errorf("consumer is closed")
	}
}

// PollMessages retrieves multiple messages from the consumer
func (c *Consumer) PollMessages(ctx context.Context, maxMessages int) ([]*ConsumerMessage, error) {
	messages := make([]*ConsumerMessage, 0, maxMessages)

	for i := 0; i < maxMessages; i++ {
		msg, err := c.Poll(ctx)
		if err != nil {
			// If we got some messages, return them
			if len(messages) > 0 {
				return messages, nil
			}
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// CommitOffset commits an offset for a topic partition
func (c *Consumer) CommitOffset(topic string, partition int32, offset int64) error {
	c.offsetsMutex.Lock()
	defer c.offsetsMutex.Unlock()

	if _, ok := c.currentOffsets[topic]; !ok {
		c.currentOffsets[topic] = make(map[int32]int64)
	}

	c.currentOffsets[topic][partition] = offset

	// TODO: Send offset commit to broker
	return nil
}

// CommitMessage commits the offset of a consumed message
func (c *Consumer) CommitMessage(msg *ConsumerMessage) error {
	return c.CommitOffset(msg.Topic, msg.Partition, msg.Offset)
}

// messageFetcher continuously fetches messages from the broker
func (c *Consumer) messageFetcher() {
	defer c.wg.Done()

	fetchTicker := time.NewTicker(
		time.Duration(c.config.FetchTimeoutMs) * time.Millisecond,
	)
	defer fetchTicker.Stop()

	for {
		select {
		case <-fetchTicker.C:
			// TODO: Fetch messages from broker
			// For now, this is a placeholder

		case <-c.ctx.Done():
			return
		}
	}
}

// heartbeatSender sends periodic heartbeats to the broker
func (c *Consumer) heartbeatSender() {
	defer c.wg.Done()

	for {
		select {
		case <-c.heartbeatTicker.C:
			// TODO: Send heartbeat to broker

		case <-c.ctx.Done():
			return
		}
	}
}

// Subscribe subscribes to additional topics
func (c *Consumer) Subscribe(topics []string) error {
	c.closedMutex.Lock()
	if c.closed {
		c.closedMutex.Unlock()
		return fmt.Errorf("consumer is closed")
	}
	c.closedMutex.Unlock()

	c.offsetsMutex.Lock()
	defer c.offsetsMutex.Unlock()

	for _, topic := range topics {
		if _, ok := c.currentOffsets[topic]; !ok {
			c.currentOffsets[topic] = make(map[int32]int64)
		}
	}

	return nil
}

// Close gracefully shuts down the consumer
func (c *Consumer) Close() error {
	c.closedMutex.Lock()
	if c.closed {
		c.closedMutex.Unlock()
		return nil
	}
	c.closed = true
	c.closedMutex.Unlock()

	c.cancel()
	c.heartbeatTicker.Stop()
	c.wg.Wait()

	if c.connection != nil {
		c.connection.Close()
	}

	close(c.messagesChan)

	return nil
}
