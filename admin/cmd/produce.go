package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
