package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// NewTopicCmd creates a new topic command
func NewTopicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "topic",
		Short: "Manage topics",
		Long:  "Create, delete, list, and describe topics",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "create <topic-name>",
		Short: "Create a new topic",
		Long:  "Create a new topic with the specified name and configuration",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return topicCreate(args[0], cmd)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "delete <topic-name>",
		Short: "Delete a topic",
		Long:  "Delete an existing topic and all its data",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return topicDelete(args[0], cmd)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all topics",
		Long:  "List all topics in the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return topicList(cmd)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "describe <topic-name>",
		Short: "Describe a topic",
		Long:  "Show detailed information about a topic",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return topicDescribe(args[0], cmd)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "alter <topic-name>",
		Short: "Alter a topic configuration",
		Long:  "Modify topic configuration like partitions and replication factor",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return topicAlter(args[0], cmd)
		},
	})

	cmd.PersistentFlags().IntP("partitions", "p", 1, "Number of partitions")
	cmd.PersistentFlags().IntP("replication-factor", "r", 1, "Replication factor")
	cmd.PersistentFlags().StringToStringP("config", "c", nil, "Topic configuration")

	return cmd
}

// NewProduceCmd creates a new produce command
func NewProduceCmd() *cobra.Command {
	var (
		topic     string
		partition int32
		key       string
		value     string
		fromFile  string
	)

	cmd := &cobra.Command{
		Use:   "produce",
		Short: "Produce messages to a topic",
		Long:  "Send messages to a StreamForge topic",
		RunE: func(cmd *cobra.Command, args []string) error {
			return produceMessages(topic, partition, key, value, fromFile)
		},
	}

	cmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic name (required)")
	cmd.Flags().Int32VarP(&partition, "partition", "p", -1, "Partition number (auto if -1)")
	cmd.Flags().StringVarP(&key, "key", "k", "", "Message key")
	cmd.Flags().StringVarP(&value, "value", "v", "", "Message value")
	cmd.Flags().StringVarP(&fromFile, "file", "f", "", "Read messages from file (one per line)")

	cmd.MarkFlagRequired("topic")

	return cmd
}

// NewConsumeCmd creates a new consume command
func NewConsumeCmd() *cobra.Command {
	var (
		topic        string
		groupID      string
		offset       string
		maxMessages  int
		timeout      int64
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

// Topic command implementations
func topicCreate(topicName string, cmd *cobra.Command) error {
	partitions, _ := cmd.Flags().GetInt("partitions")
	replicationFactor, _ := cmd.Flags().GetInt("replication-factor")
	config, _ := cmd.Flags().GetStringToString("config")

	fmt.Printf("Creating topic '%s'\n", topicName)
	fmt.Printf("  Partitions: %d\n", partitions)
	fmt.Printf("  Replication Factor: %d\n", replicationFactor)
	if len(config) > 0 {
		fmt.Printf("  Config: %v\n", config)
	}

	// TODO: Call broker API to create topic
	fmt.Println("✓ Topic created successfully")
	return nil
}

func topicDelete(topicName string, cmd *cobra.Command) error {
	fmt.Printf("Deleting topic '%s'\n", topicName)

	// TODO: Call broker API to delete topic
	fmt.Println("✓ Topic deleted successfully")
	return nil
}

func topicList(cmd *cobra.Command) error {
	fmt.Println("Topics in cluster:")
	fmt.Println("-----------------")

	// TODO: Call broker API to list topics
	topics := []string{"events", "logs", "metrics"}
	for i, topic := range topics {
		fmt.Printf("%d. %s\n", i+1, topic)
	}

	return nil
}

func topicDescribe(topicName string, cmd *cobra.Command) error {
	fmt.Printf("Topic: %s\n", topicName)
	fmt.Println("Partitions:")
	fmt.Println("-----------")

	// TODO: Call broker API to describe topic
	fmt.Printf("  Partition 0 | Leader: 1 | Replicas: [1,2,3] | ISR: [1,2,3]\n")
	fmt.Printf("  Partition 1 | Leader: 2 | Replicas: [2,3,1] | ISR: [2,3,1]\n")
	fmt.Printf("  Partition 2 | Leader: 3 | Replicas: [3,1,2] | ISR: [3,1,2]\n")

	return nil
}

func topicAlter(topicName string, cmd *cobra.Command) error {
	partitions, _ := cmd.Flags().GetInt("partitions")

	fmt.Printf("Altering topic '%s'\n", topicName)
	if partitions > 0 {
		fmt.Printf("  New Partitions: %d\n", partitions)
	}

	// TODO: Call broker API to alter topic
	fmt.Println("✓ Topic altered successfully")
	return nil
}

// Produce command implementations
func produceMessages(topic string, partition int32, key string, value string, fromFile string) error {
	if fromFile != "" {
		return produceFromFile(topic, partition, fromFile)
	}

	if value == "" {
		return fmt.Errorf("either --value or --file must be specified")
	}

	fmt.Printf("Producing to topic: %s\n", topic)
	if partition >= 0 {
		fmt.Printf("  Partition: %d\n", partition)
	} else {
		fmt.Println("  Partition: auto-selected")
	}
	if key != "" {
		fmt.Printf("  Key: %s\n", key)
	}
	fmt.Printf("  Value: %s\n", value)

	// TODO: Call producer API
	fmt.Println("✓ Message published successfully")
	return nil
}

func produceFromFile(topic string, partition int32, filePath string) error {
	fmt.Printf("Reading messages from file: %s\n", filePath)

	// TODO: Read file and produce each line as a message
	// For now, simulate reading the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fmt.Printf("Producing %d messages to topic '%s'\n", 5, topic) // Placeholder
	fmt.Println("✓ All messages published successfully")

	return nil
}

// Consume command implementations
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
			"offset":    0,
			"key":       "key1",
			"value":     "message 1",
			"timestamp": time.Now(),
		},
		{
			"offset":    1,
			"key":       "key2",
			"value":     "message 2",
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
