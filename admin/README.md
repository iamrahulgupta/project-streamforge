# StreamForge Admin CLI & REST API

A comprehensive administration tool for StreamForge with both CLI and REST API interfaces for managing brokers, topics, and consumer groups.

## Features

- **CLI Tool**: Command-line interface for managing StreamForge
- **REST API**: RESTful API server for remote administration
- **Topic Management**: Create, delete, list, and describe topics
- **Message Operations**: Produce and consume messages from the CLI
- **Consumer Groups**: List and describe consumer groups
- **Cluster Information**: View broker and cluster metadata
- **Multiple Output Formats**: JSON, CSV, and text output options

## Installation

```bash
cd admin
go build -o streamforge-admin .
```

Or for development:

```bash
go run main.go --help
```

## CLI Usage

### Topic Management

```bash
# List all topics
./streamforge-admin topic list

# Create a topic
./streamforge-admin topic create events \
  --partitions 3 \
  --replication-factor 3

# Describe a topic
./streamforge-admin topic describe events

# Alter topic configuration
./streamforge-admin topic alter events --partitions 5

# Delete a topic
./streamforge-admin topic delete events
```

### Message Produce

```bash
# Produce a single message
./streamforge-admin produce \
  --topic events \
  --key user123 \
  --value "user login event"

# Produce messages from file (one message per line)
./streamforge-admin produce \
  --topic events \
  --file messages.txt

# Produce to specific partition
./streamforge-admin produce \
  --topic events \
  --partition 0 \
  --value "important event"
```

### Message Consume

```bash
# Consume messages from latest offset
./streamforge-admin consume \
  --topic events \
  --group my-consumer-group

# Consume from earliest offset
./streamforge-admin consume \
  --topic events \
  --group my-consumer-group \
  --offset earliest

# Consume with JSON output
./streamforge-admin consume \
  --topic events \
  --group my-consumer-group \
  --format json \
  --max-messages 10

# Consume with CSV output
./streamforge-admin consume \
  --topic events \
  --group my-consumer-group \
  --format csv

# Consume from specific offset
./streamforge-admin consume \
  --topic events \
  --group my-consumer-group \
  --offset 100
```

### Broker Information

```bash
# Connect to specific brokers
./streamforge-admin --brokers localhost:9092,localhost:9093 \
  topic list

# Set brokers via environment variable
export STREAMFORGE_BROKERS=localhost:9092,localhost:9093
./streamforge-admin topic list
```

## REST API Usage

### Start the Server

```bash
# Default port 8080
./streamforge-admin server

# Custom port
./streamforge-admin server --port 8888

# With specific brokers
./streamforge-admin --brokers localhost:9092,localhost:9093 server --port 8080
```

### API Endpoints

#### Health Check
```bash
GET /health

Example:
curl http://localhost:8080/health
```

#### Brokers
```bash
GET /api/v1/brokers
```

#### Topics
```bash
# List topics
GET /api/v1/topics

# Create topic
POST /api/v1/topics
Content-Type: application/json
{
  "name": "events",
  "partitions": 3,
  "replication_factor": 3,
  "config": {}
}

# Get topic details
GET /api/v1/topics/{topic-name}

# Delete topic
DELETE /api/v1/topics/{topic-name}
```

#### Consumer Groups
```bash
# List groups
GET /api/v1/groups

# Get group details
GET /api/v1/groups/{group-id}
```

#### Produce Messages
```bash
POST /api/v1/produce
Content-Type: application/json
{
  "topic": "events",
  "key": "user123",
  "value": "message content"
}
```

#### Consume Messages
```bash
POST /api/v1/consume
Content-Type: application/json
{
  "topic": "events",
  "group_id": "my-group",
  "max_messages": 10,
  "offset": "latest"
}
```

## REST API Examples

### List all topics
```bash
curl http://localhost:8080/api/v1/topics
```

### Create a topic
```bash
curl -X POST http://localhost:8080/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{
    "name": "orders",
    "partitions": 5,
    "replication_factor": 3
  }'
```

### Get topic details
```bash
curl http://localhost:8080/api/v1/topics/events
```

### Produce a message
```bash
curl -X POST http://localhost:8080/api/v1/produce \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "events",
    "key": "order-123",
    "value": "Order placed successfully"
  }'
```

### Consume messages
```bash
curl -X POST http://localhost:8080/api/v1/consume \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "events",
    "group_id": "order-processors",
    "max_messages": 5,
    "offset": "latest"
  }'
```

## Project Structure

```
admin/
├── main.go                  # Entry point and CLI setup
├── server/
│   └── rest.go             # REST API server implementation
├── cmd/
│   ├── root.go             # Helper functions and shared code
│   ├── topic.go            # Topic management commands
│   ├── produce.go          # Message produce commands
│   └── consume.go          # Message consume commands
└── README.md               # This file
```

## Configuration

### Environment Variables

```bash
# Broker addresses (comma-separated)
export STREAMFORGE_BROKERS=localhost:9092,localhost:9093

# Default consumer group
export STREAMFORGE_GROUP=default-group

# REST API port
export STREAMFORGE_API_PORT=8080
```

### Command Flags

All commands support these global flags:

```
--brokers string[]    Broker addresses (default: localhost:9092)
--help                Show help
--version             Show version
```

## Build

```bash
# Build binary
go build -o streamforge-admin .

# Build with version info
go build -ldflags "-X main.version=0.1.0" -o streamforge-admin .

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o streamforge-admin-linux .

# Cross-compile for macOS
GOOS=darwin GOARCH=amd64 go build -o streamforge-admin-darwin .
```

## Testing

```bash
# Run tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestTopicCreate ./cmd
```

## Key Commands Reference

```
streamforge --help                           Show general help
streamforge topic --help                     Show topic command help
streamforge produce --help                   Show produce command help
streamforge consume --help                   Show consume command help
streamforge server --help                    Show server command help

streamforge topic list                       List all topics
streamforge topic create <name>              Create topic
streamforge topic describe <name>            Get topic details
streamforge topic delete <name>              Delete topic

streamforge produce -t <topic> -v <value>    Produce message
streamforge consume -t <topic> -g <group>    Consume messages

streamforge server --port 8080               Start REST API server
```

## Performance Tips

- Use `--max-messages` to limit consumption for large topics
- Batch produce operations for better throughput
- Use compression for large messages (via topic config)
- Connection pooling is handled automatically for REST API

## Troubleshooting

### Connection Errors
```bash
# Verify broker is running
./streamforge-admin --brokers localhost:9092 topic list

# Try alternative broker
./streamforge-admin --brokers 192.168.1.100:9092 topic list
```

### Consumer Group Issues
```bash
# Verify group exists
./streamforge-admin groups

# Check group status
./streamforge-admin group describe my-group
```

## Contributing

Contributions are welcome! Please follow the development guidelines in the main README.

## License

This project is licensed under the MIT License.
