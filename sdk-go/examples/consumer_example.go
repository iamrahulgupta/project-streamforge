package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iamrahulgupta/streamforge-go/consumer"
)

// Example of using the StreamForge Consumer SDK
func main() {
	brokers := []string{"localhost:9092"}
	groupID := "my-consumer-group"
	topics := []string{"events"}

	// Create consumer with default config
	config := consumer.DefaultConsumerConfig(brokers, groupID, topics)
	c, err := consumer.NewConsumer(config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()

	fmt.Printf("Consumer created for group '%s' subscribing to %v\n", groupID, topics)

	// Example 1: Poll messages one at a time
	fmt.Println("\nPolling messages (with timeout)...")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	for i := 0; i < 3; i++ {
		msg, err := c.Poll(ctx)
		if err != nil {
			fmt.Printf("Poll timeout or error: %v\n", err)
			break
		}

		fmt.Printf("Received message: topic=%s, partition=%d, offset=%d\n",
			msg.Topic, msg.Partition, msg.Offset)
		fmt.Printf("  Key: %s, Value: %s\n", string(msg.Key), string(msg.Value))

		// Commit the offset after processing
		if err := c.CommitMessage(msg); err != nil {
			log.Printf("Error committing offset: %v", err)
		}
	}

	// Example 2: Poll multiple messages
	fmt.Println("\nPolling batch of messages...")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	messages, err := c.PollMessages(ctx, 5)
	if err != nil {
		fmt.Printf("Error polling messages: %v\n", err)
	} else {
		fmt.Printf("Received %d messages\n", len(messages))
		for _, msg := range messages {
			fmt.Printf("  Message: offset=%d, key=%s\n",
				msg.Offset, string(msg.Key))

			if err := c.CommitMessage(msg); err != nil {
				log.Printf("Error committing: %v", err)
			}
		}
	}

	// Example 3: ProcessState messages in a loop with rebalancing
	fmt.Println("\nContinuous consumption (will stop after timeout)...")
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messageCount := 0
	for {
		msg, err := c.Poll(ctx)
		if err != nil {
			if err == context.DeadlineExceeded {
				fmt.Println("Consumer stopped (timeout)")
			} else {
				fmt.Printf("Error: %v\n", err)
			}
			break
		}

		messageCount++
		fmt.Printf("Message %d: offset=%d, value=%s\n",
			messageCount, msg.Offset, string(msg.Value))

		// Commit offset
		c.CommitMessage(msg)

		// Process message (simulate work)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\nConsumed %d total messages\n", messageCount)
	fmt.Println("Consumer example completed")
}
