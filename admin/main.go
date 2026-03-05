package main

import (
	"fmt"
	"os"

	"github.com/iamrahulgupta/streamforge-admin/server"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	brokers []string
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "streamforge",
		Short:   "StreamForge - Distributed Event Streaming Platform",
		Long:    "A command-line tool for managing StreamForge brokers, topics, and consumer groups",
		Version: version,
	}

	// Global flags
	rootCmd.PersistentFlags().StringSliceVar(
		&brokers,
		"brokers",
		[]string{"localhost:9092"},
		"Comma-separated list of broker addresses",
	)

	// Add subcommands
	rootCmd.AddCommand(NewTopicCmd())
	rootCmd.AddCommand(NewProduceCmd())
	rootCmd.AddCommand(NewConsumeCmd())
	rootCmd.AddCommand(NewServerCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

// NewServerCmd returns the server command
func NewServerCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the admin REST API server",
		Long:  "Start the StreamForge admin REST API server for remote management",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startServer(port, brokers)
		},
	}

	cmd.Flags().IntVar(&port, "port", 8080, "Port to listen on")

	return cmd
}

func startServer(port int, brokerAddrs []string) error {
	s := server.NewServer(port, brokerAddrs)
	return s.Start()
}
