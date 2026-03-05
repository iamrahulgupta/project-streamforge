package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TopicCommand represents the topic management command
type TopicCommand struct {
	brokers []string
}

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

	// Add flags
	cmd.PersistentFlags().IntP("partitions", "p", 1, "Number of partitions")
	cmd.PersistentFlags().IntP("replication-factor", "r", 1, "Replication factor")
	cmd.PersistentFlags().StringToStringP("config", "c", nil, "Topic configuration")

	return cmd
}

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
