package main

import (
	"context"
	"fmt"
	"log"

	"github.com/iamrahulgupta/streamforge-go/producer"
)

// Example of using the StreamForge Producer SDK
func main() {
	brokers := []string{"localhost:9092"}
	topic := "events"

	// Create producer with default config
	config := producer.DefaultProducerConfig(brokers, topic)
	p, err := producer.NewProducer(config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()

	// Example 1: Async produce (fire and forget)
	fmt.Println("Producing messages asynchronously...")
	for i := 0; i < 5; i++ {
		msg := &producer.Message{
			Key:   []byte(fmt.Sprintf("key-%d", i)),
			Value: []byte(fmt.Sprintf("message-%d", i)),
			Headers: map[string]string{
				"source": "example-producer",
			},
		}

		_, err := p.Produce(context.Background(), msg)
		if err != nil {
			log.Printf("Error producing message: %v", err)
		}
	}

	// Example 2: Sync produce (wait for confirmation)
	fmt.Println("\nProducing message with sync confirmation...")
	msg := &producer.Message{
		Key:   []byte("important-key"),
		Value: []byte("important message"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5_000) // 5 seconds timeout
	result, err := p.ProduceSync(ctx, msg)
	cancel()

	if err != nil {
		log.Fatalf("Error producing message: %v", err)
	}

	fmt.Printf("Message produced: Topic=%s, Partition=%d, Offset=%d\n",
		result.Topic, result.Partition, result.Offset)

	// Example 3: Produce multiple messages
	fmt.Println("\nProducing batch of messages...")
	for i := 0; i < 10; i++ {
		msg := &producer.Message{
			Key:   []byte(fmt.Sprintf("batch-key-%d", i)),
			Value: []byte(fmt.Sprintf("batch-message-%d with some data", i)),
		}

		p.Produce(context.Background(), msg)
	}

	// Wait for all messages to be sent
	if err := p.Flush(10_000); err != nil {
		log.Printf("Error flushing: %v", err)
	}

	fmt.Println("\nAll messages sent successfully!")
}
