package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// NewConsumeCmd creates a new consume command
func NewConsumeCmd() *cobra.Command {
	var (
		topic      string
		groupID    string
		offset     string
		maxMessages int
		timeout    int64
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "consume",
		Short: "Consume messages from a topic",
		Long:  "Read messages from a StreamForge topic",
		RunE: func(cmd *cobra.Command, args []string) error {
			return consumeMessages(topic, groupID, offset, maxMessages, timeout, outputFormat)
		},
	}

	cmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic name (required)")
	cmd.Flags().StringVarP(&groupID, "group", "g", "", "Consumer group ID (required)")
	cmd.Flags().StringVarP(&offset, "offset", "o", "latest", "Starting offset: earliest, latest, or offset number")
	cmd.Flags().IntVarP(&maxMessages, "max-messages", "m", 10, "Maximum number of messages to consume")
	cmd.Flags().Int64VarP(&timeout, "timeout", "T", 30000, "Timeout in milliseconds")
	cmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format: json, csv, text")

	cmd.MarkFlagRequired("topic")
	cmd.MarkFlagRequired("group")

	return cmd
}

func consumeMessages(topic string, groupID string, offset string, maxMessages int, timeoutMs int64, format string) error {
	fmt.Printf("Consuming from topic: %s\n", topic)
	fmt.Printf("  Group ID: %s\n", groupID)
	fmt.Printf("  Starting Offset: %s\n", offset)
	fmt.Printf("  Max Messages: %d\n", maxMessages)
	fmt.Printf("  Timeout: %dms\n", timeoutMs)
	fmt.Printf("  Format: %s\n", format)
	fmt.Println()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(timeoutMs)*time.Millisecond,
	)
	defer cancel()

	// TODO: Call consumer API
	// Simulate consuming messages
	messages := []map[string]interface{}{
		{
			"offset": 0,
			"key":    "key1",
			"value":  "message 1",
			"timestamp": time.Now(),
		},
		{
			"offset": 1,
			"key":    "key2",
			"value":  "message 2",
			"timestamp": time.Now(),
		},
	}

	fmt.Println("Messages:")
	fmt.Println("---------")

	for i, msg := range messages {
		if i >= maxMessages {
			break
		}

		switch format {
		case "json":
			fmt.Printf("{\n")
			fmt.Printf("  \"offset\": %v,\n", msg["offset"])
			fmt.Printf("  \"key\": \"%v\",\n", msg["key"])
			fmt.Printf("  \"value\": \"%v\",\n", msg["value"])
			fmt.Printf("  \"timestamp\": \"%v\"\n", msg["timestamp"])
			fmt.Printf("}\n")
		case "csv":
			fmt.Printf("%v,%v,%v,%v\n", msg["offset"], msg["key"], msg["value"], msg["timestamp"])
		default: // text
			fmt.Printf("Offset: %v | Key: %v | Value: %v | Timestamp: %v\n",
				msg["offset"], msg["key"], msg["value"], msg["timestamp"])
		}
	}

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("\n✓ Consume timeout reached")
		}
	default:
		fmt.Println("\n✓ Consume completed")
	}

	return nil
}
