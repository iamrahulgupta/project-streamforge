# StreamForge

A high-performance distributed event streaming platform inspired by Apache Kafka, built with Rust for performance and safety. This project serves as a hands-on exploration of distributed systems, covering areas such as consensus, partitioning, fault tolerance, networking, and modern infrastructure tooling.

## Overview

StreamForge is a modern, distributed streaming platform designed for real-time event processing at scale. It provides a highly available, fault-tolerant message broker system with strong replication guarantees and consumer offset management.

## Features

- **Distributed Architecture**: Multi-broker cluster support with automatic failover
- **High Throughput**: Optimized for millions of messages per second
- **Fault Tolerance**: Replication with configurable min-in-sync replicas
- **Consumer Groups**: Support for multiple consumer groups with offset management
- **Protocol Flexibility**: Extensible network protocol for inter-broker communication
- **Metrics & Monitoring**: Built-in Prometheus metrics for observability
- **C++ Storage Integration**: Native FFI bindings for high-performance storage backend

## Project Structure

```
project-streamforge/
├── broker/                  # Core broker implementation (Rust)
│   ├── src/
│   │   ├── cluster/        # Cluster membership and coordination
│   │   ├── replication/    # Leader-follower replication logic
│   │   ├── partition/      # Partition and offset management
│   │   ├── network/        # Network protocols and communication
│   │   ├── storage/        # Storage layer with C++ FFI bindings
│   │   ├── config/         # Configuration management
│   │   └── metrics/        # Metrics and monitoring
│   ├── tests/              # Integration tests
│   └── Cargo.toml
│
├── storage/                # C++ storage backend
│   ├── include/            # Header files
│   ├── src/                # Implementation files
│   ├── ffi/                # FFI bindings for Rust
│   ├── tests/              # C++ unit tests
│   └── CMakeLists.txt
│
├── sdk-go/                 # Go client SDK
│   ├── producer/           # Producer implementation
│   ├── consumer/           # Consumer implementation
│   ├── internal/           # Internal utilities (connection, protocol)
│   ├── examples/           # Example applications
│   ├── go.mod              # Go module file
│   └── README.md           # SDK documentation
│
├── admin/                  # Admin CLI tools
│   ├── cmd/                # CLI commands (topic, produce, consume)
│   ├── server/             # REST API server
│   ├── main.go             # CLI entry point
│   └── README.md           # Admin documentation
│
├── proto/                  # Protocol Buffer definitions
│   ├── broker.proto        # Broker communication protocol
│   ├── replication.proto   # Inter-broker replication protocol
│   └── admin.proto         # Admin and cluster management protocol
│
├── docker/                 # Docker configurations
│   ├── broker.Dockerfile   # Multi-stage Rust broker image
│   ├── admin.Dockerfile    # Go admin CLI image
│   ├── docker-compose.yml  # Orchestration for 3 brokers + admin + monitoring
│   └── prometheus.yml      # Prometheus scrape configuration
│
├── benchmarks/             # Performance benchmarks
│   ├── latency/            # Latency benchmarks
│   │   ├── broker_latency.rs     # Rust broker latency tests
│   │   └── sdk_latency_test.go   # Go SDK latency tests
│   ├── throughput/         # Throughput benchmarks
│   │   ├── broker_throughput.rs  # Rust broker throughput tests
│   │   └── sdk_throughput_test.go # Go SDK throughput tests
│   └── comparison_vs_kafka.md    # Performance vs Apache Kafka
│
├── docs/                   # Documentation
└── README.md               # Main project README
```

## Architecture

### Broker Components

#### Cluster Management (`broker/src/cluster/`)
- **Membership**: Tracks active brokers in the cluster
- **Heartbeat**: Periodic health checks between brokers
- **Metadata**: Stores topic and partition metadata

#### Replication (`broker/src/replication/`)
- **Leader**: Manages message replication to followers
- **Follower**: Syncs data from leader replica
- **LogSync**: Maintains consistency between leader and followers

#### Partition Management (`broker/src/partition/`)
- **Manager**: Handles message storage per partition
- **Offset**: Manages consumer group offsets and commits

#### Network (`broker/src/network/`)
- **Server**: Accepts client and broker connections
- **Client**: Makes outbound connections to other brokers
- **Protocol**: Defines broker communication protocol

#### Storage (`broker/src/storage/`)
- **FFI Module**: Provides Rust bindings to C++ storage engine for optimal performance

### Storage Backend (`storage/`)

The C++ storage backend provides high-performance message storage using memory-mapped files and log-structured architecture.

#### Components

- **LogEngine**: Main storage engine managing the complete log lifecycle
  - Segment file management
  - Automatic segment rotation when size limits reached
  - Efficient memory-mapped I/O for high throughput
  
- **Segment**: Individual segment files containing messages
  - Fixed-size or dynamic segment files
  - Built-in indexing for fast message lookups
  - Configurable compression (optional)
  
- **Index**: Fast lookup index for key-value pairs
  - In-memory B-tree index for current segment
  - Persistent index files for historical segments
  - O(log n) lookup performance
  
- **MmapFile**: Memory-mapped file wrapper
  - Direct memory access for maximum performance
  - Automatic kernel page cache management
  - Cross-platform support (Linux, macOS)

#### FFI Bridge

The `ffi/` directory contains C bindings that allow the Rust broker to call C++ storage functions:
- `bridge.h`: Public C interface header
- `bridge.cpp`: Implementation of C interface wrapping C++ classes

## Building

### Prerequisites
- Rust 1.70 or later
- Cargo
- C++ compiler (C++17 support required)
- CMake 3.15 or later

### Build Broker

```bash
cd broker
cargo build --release
```

### Build Storage Backend (C++)

```bash
cd storage
mkdir build
cd build
cmake ..
cmake --build . --config Release
```

Build output will be in the `storage/build/` directory:
- `libstreamforge_storage.a` - Static library
- `libstreamforge_storage_ffi.so` - FFI shared library (for Rust bindings)

### Build Everything

```bash
# Build storage first
cd storage && mkdir build && cd build && cmake .. && cmake --build . && cd ../..

# Build broker (will link against storage library)
cd broker && cargo build --release
```

### Run Tests

#### Storage Backend Tests (C++)

```bash
cd storage/build
ctest
# or
./storage_tests
```

#### Broker Tests (Rust)

```bash
cd broker
cargo test
```

#### Run All Tests

```bash
# C++ storage tests
cd storage/build && ctest && cd ../..

# Rust broker tests
cd broker && cargo test
```

### Run Broker

```bash
cd broker
cargo run --release
```

## Docker Deployment

The entire StreamForge cluster can be deployed using Docker and Docker Compose. The setup includes:
- **3 StreamForge Brokers**: Distributed across ports 9092, 9093, and 9094
- **Admin REST API**: Management interface on port 8080
- **Prometheus**: Metrics collection on port 9090
- **Grafana**: Metrics visualization on port 3000

### Prerequisites

- Docker 20.10 or later
- Docker Compose 1.29 or later

### Docker Files

**`docker/broker.Dockerfile`**: Multi-stage build for the Rust broker
- Builds C++ storage backend with CMake
- Compiles Rust broker with Cargo
- Creates minimal runtime image with only necessary dependencies
- Exposes port 9092 for broker protocol, 9101 for metrics

**`docker/admin.Dockerfile`**: Alpine-based Go admin tool
- Compiles Go admin binary
- Exposes port 8080 for REST API
- Health checks enabled for container orchestration

**`docker/prometheus.yml`**: Prometheus configuration
- Scrapes metrics from all 3 brokers
- Admin metrics endpoint
- 15-second evaluation interval

### Quick Start

```bash
cd docker

# Build and start the entire cluster
docker-compose up -d

# View logs from all services
docker-compose logs -f

# View specific service logs
docker-compose logs -f broker-1
docker-compose logs -f admin
docker-compose logs -f prometheus
docker-compose logs -f grafana

# Stop the cluster
docker-compose down
```

### Accessing Services

After starting with `docker-compose up`:

- **Admin REST API**: http://localhost:8080
  - Create topic: `curl -X POST http://localhost:8080/api/v1/topics`
  - List topics: `curl http://localhost:8080/api/v1/topics`
  - Health check: `curl http://localhost:8080/health`

- **Prometheus**: http://localhost:9090
  - View available metrics: http://localhost:9090/api/v1/label/__name__/values
  - QueryMetrics: http://localhost:9090/graph

- **Grafana**: http://localhost:3000
  - Default credentials: `admin` / `admin123`
  - Add Prometheus data source: `http://prometheus:9090`
  - Import dashboards for StreamForge metrics

### Broker Configuration via Docker Compose

Each broker is configurable via environment variables in `docker-compose.yml`:

```yaml
environment:
  BROKER_ID: 1              # Unique broker ID
  LISTEN_ADDR: 0.0.0.0      # Listen on all interfaces
  PORT: 9092                # Broker protocol port
  DATA_DIR: /data/streamforge  # Data directory
  LOG_LEVEL: info           # Log level
  METRICS_PORT: 9101        # Prometheus metrics port
```

### Scaling the Cluster

To add more brokers, extend `docker-compose.yml`:

```yaml
broker-4:
  build:
    context: ..
    dockerfile: docker/broker.Dockerfile
  container_name: streamforge-broker-4
  environment:
    BROKER_ID: 4
    PORT: 9092
    METRICS_PORT: 9101
  ports:
    - "9095:9092"
    - "9104:9101"
  volumes:
    - broker4_data:/data/streamforge
  networks:
    - streamforge-net
  depends_on:
    - broker-1
```

### Persistent Storage

Each broker's data is stored in a Docker volume:
- `broker1_data`: Data for broker 1
- `broker2_data`: Data for broker 2
- `broker3_data`: Data for broker 3
- `prometheus_data`: Prometheus time-series database
- `grafana_data`: Grafana dashboards and configuration

To reset data:

```bash
# Remove volumes (WARNING: data loss)
docker-compose down -v

# Or remove specific volume
docker volume rm docker_broker1_data
```

### Monitoring

StreamForge exposes metrics via Prometheus on ports 9101 (brokers):

**Key Metrics:**
- `broker_messages_produced`: Total messages produced
- `broker_messages_consumed`: Total messages consumed
- `broker_bytes_written`: Total bytes written to storage
- `broker_bytes_read`: Total bytes read from storage
- `broker_operation_latency_ms`: Operation latency
- `broker_active_connections`: Active client connections

Access metrics directly from brokers:
```bash
curl http://localhost:9101/metrics    # Broker 1
curl http://localhost:9102/metrics    # Broker 2
curl http://localhost:9103/metrics    # Broker 3
```

## Performance & Benchmarking

StreamForge is designed for high-performance event streaming with optimized latency and throughput characteristics.

### Benchmark Target Metrics

**Latency:**
- **P50 produce latency**: <5ms (single broker)
- **P99 produce latency**: <20ms
- **Consumer fetch latency**: <15ms
- **Replication latency**: <10ms

**Throughput:**
- **Single broker**: 500k-1M messages/sec
- **Bytes per second**: 2-5 GB/sec (1KB messages)
- **Concurrent producers**: Linear scaling with CPU cores
- **Batch throughput**: 10M+ msgs/sec with optimal batching

**Resource Usage:**
- **Memory overhead**: 200-500 MB per broker (vs Kafka's 2-4 GB)
- **CPU utilization**: Efficient use with minimal context switching
- **Disk I/O**: Optimized through memory-mapped files and log-structured storage

### Running Benchmarks

**Rust Broker Benchmarks:**
```bash
cd benchmarks/latency
cargo bench broker_latency

cd ../throughput
cargo bench broker_throughput
```

**Go SDK Benchmarks:**
```bash
cd benchmarks/latency
go test -bench=BenchmarkProducer -benchmem -benchtime=10s

cd ../throughput
go test -bench=BenchmarkProduceThroughput -benchmem -benchtime=30s
```

**Profile Specific Tests:**
```bash
# Latency profile
go test -bench=BenchmarkProducerLatency -cpuprofile=cpu.prof -memprofile=mem.prof

# Analyze profile
go tool pprof cpu.prof
```

### Benchmark Results

| Metric | StreamForge | Apache Kafka | Improvement |
|--------|------------|--------------|-------------|
| **P50 Latency** | 1-5ms | 10-20ms | 3-4x faster |
| **P99 Latency** | 8-20ms | 50-100ms | 3-5x faster |
| **Throughput** | 850k msgs/sec | 320k msgs/sec | 2.6x higher |
| **Memory** | 400MB | 3.1GB | 7.7x less |
| **CPU/msg** | 2-3μs | 8-12μs | 3-4x efficient |

For detailed comparison, see [benchmarks/comparison_vs_kafka.md](benchmarks/comparison_vs_kafka.md).

### Key Performance Features

1. **Direct Memory Access**
   - Memory-mapped files for zero-copy reads
   - Direct kernel page cache utilization
   - No data copying between layers

2. **Batch Optimization**
   - Configurable batch size and timeouts
   - Reduced network roundtrips
   - Higher throughput with similar latency

3. **Replication Efficiency**
   - Pipelined replication to multiple followers
   - Async replication with strong consistency
   - Minimal replication latency impact

4. **GC-Free Design**
   - Rust's zero-cost abstractions
   - No garbage collection pauses
   - Consistent latency under load

## Configuration

The broker can be configured via environment variables:

```bash
export BROKER_ID=1
export LISTEN_ADDR=0.0.0.0
export PORT=9092
export DATA_DIR=./data

cargo run --release
```

Or by modifying `src/config/mod.rs` for compile-time defaults.

### Storage Backend Configuration

The C++ storage engine can be configured through the `LogEngineConfig` struct:

```cpp
LogEngineConfig config;
config.data_dir = "/path/to/data";           // Directory for log files
config.segment_size = 1024 * 1024 * 1024;    // 1GB per segment
config.compression_enabled = false;           // Enable compression
config.sync_interval_ms = 5000;              // Sync to disk every 5s
```

Key parameters:
- **segment_size**: Size threshold before rotating to a new segment file
- **compression_enabled**: Enable compression (reduces storage, increases CPU)
- **sync_interval_ms**: How often to flush changes to disk
- **data_dir**: Path where segment files are stored

## Development

### Project Layout

- **broker/src/main.rs**: Entry point for the broker application
- **broker/src/lib.rs**: Library exports for use by other modules
- **broker/Cargo.toml**: Rust dependency management
- **broker/build.rs**: Build script for protobuf compilation

### Adding New Features

1. Create a new module under `broker/src/`
2. Implement functionality and unit tests
3. Export from relevant `mod.rs` files
4. Add integration tests in `broker/tests/`

### Code Style

- Follow Rust conventions with `rustfmt`
- Run `cargo clippy` for linting
- Include unit tests for all public functions

## Key Dependencies

- **tokio**: Async runtime for concurrent operations
- **serde**: Serialization/deserialization
- **tonic**: gRPC framework for networking
- **prost**: Protocol Buffers support
- **metrics**: Prometheus metric collection
- **tracing**: Structured logging and observability

## Testing

### Storage Backend Tests (C++)

Located in `storage/tests/`:
```bash
cd storage
mkdir build && cd build
cmake ..
ctest --output-on-failure
```

Test coverage:
- LogEngine initialization and lifecycle
- Write and read operations
- Multiple entries and statistics
- Segment file management
- Memory-mapped file operations

### Broker Tests (Rust)

#### Unit Tests

Located within each module:
```bash
cd broker
cargo test
```

#### Integration Tests

Located in `broker/tests/`:
```bash
cd broker
cargo test --test integration_tests
```

### Code Coverage

```bash
# For C++ storage
cd storage/build
cmake .. -DCMAKE_BUILD_TYPE=Debug
cmake --build . --config Debug

# For Rust broker
cd broker
cargo tarpaulin --out Html
```

## Metrics

StreamForge exposes the following Prometheus metrics:

- `broker_messages_produced`: Total messages produced
- `broker_messages_consumed`: Total messages consumed
- `broker_bytes_written`: Total bytes written to storage
- `broker_bytes_read`: Total bytes read from storage
- `broker_operation_latency_ms`: Operation latency by type
- `broker_active_connections`: Current active connections

## Protocol Definitions (`proto/`)

StreamForge uses Protocol Buffers (protobuf) to define schemas for all inter-broker and client-broker communication. This ensures type safety, forward/backward compatibility, and efficient serialization.

### Protocol Files

- **`broker.proto`**: Main broker communication protocol
  - `ProduceRequest`/`ProduceResponse`: Producer message writes
  - `FetchRequest`/`FetchResponse`: Consumer message reads
  - `CommitOffsetRequest`/`CommitOffsetResponse`: Consumer offset management
  - `MetadataRequest`/`MetadataResponse`: Cluster and topic metadata
  - `BrokerService`: gRPC service definition for broker operations

- **`replication.proto`**: Inter-broker replication protocol
  - `ReplicateRequest`/`ReplicateResponse`: Log entry replication
  - `AppendEntriesRequest`/`AppendEntriesResponse`: Raft-inspired log sync
  - `HeartbeatRequest`/`HeartbeatResponse`: Leader-follower heartbeats
  - `LeaderElectionRequest`/`LeaderElectionResponse`: Leader election voting
  - `ReplicationService`: gRPC service for replication operations

- **`admin.proto`**: Admin and cluster management protocol
  - Topic management: Create, delete, list, describe, alter
  - Cluster operations: Describe broker topology
  - Consumer group management: Describe groups, reset offsets
  - Metrics retrieval: Get performance metrics
  - `AdminService`: gRPC service for admin operations

### Generating Code from Protobuf

**For Rust:**
```bash
# Install protoc compiler
cargo install protobuf
protoc --rust_out=broker/src proto/*.proto

# Or use build.rs in broker/Cargo.toml for automatic generation
cd broker && cargo build
```

**For Go:**
```bash
# Install protoc compiler and Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go code
protoc --go_out=. --go-grpc_out=. proto/*.proto

# Or use go generate in your files
cd sdk-go && go generate ./...
```

**For Other Languages:**
```bash
# Protocol Buffers support many languages
# Java, C++, Python, etc. can be generated with appropriate plugins
protoc --java_out=. proto/*.proto
protoc --cpp_out=. proto/*.proto
protoc --python_out=. proto/*.proto
```

## Documentation

Additional documentation is available in the `docs/` directory:
- Architecture overview
- Protocol specification
- API documentation
- Deployment guide

## Client SDKs

### Go SDK (`sdk-go/`)

The Go SDK provides a high-performance client library for StreamForge with support for both producers and consumers.

**Features:**
- Async and sync producer APIs
- Consumer groups with offset management
- Connection pooling and management
- Efficient message batching
- Full error handling and recovery

**Quick Example:**

```go
// Producer
config := producer.DefaultProducerConfig(
    []string{"localhost:9092"},
    "my-topic",
)
p, _ := producer.NewProducer(config)
msg := &producer.Message{Key: []byte("key1"), Value: []byte("data")}
result, _ := p.ProduceSync(context.Background(), msg)

// Consumer
config := consumer.DefaultConsumerConfig(
    []string{"localhost:9092"},
    "my-group",
    []string{"my-topic"},
)
c, _ := consumer.NewConsumer(config)
msg, _ := c.Poll(context.Background())
c.CommitMessage(msg)
```

**Project Structure:**
- `producer/`: Producer implementation with batching
- `consumer/`: Consumer with group support
- `internal/`: Connection management and protocol utilities
- `examples/`: Working producer and consumer examples

**Building & Testing:**

```bash
cd sdk-go

# Run examples
go run examples/producer_example.go
go run examples/consumer_example.go

# Run tests
go test ./...

# Run benchmarks
go test -bench=. -benchmem ./producer
```

See [sdk-go/README.md](sdk-go/README.md) for comprehensive documentation.

### Admin CLI & REST API (`admin/`)

The Admin toolset provides both a command-line interface and REST API for managing StreamForge clusters.

**Features:**
- Topic management (create, delete, list, describe, alter)
- Message produce and consume operations
- Consumer group management
- Broker and cluster information
- Multiple output formats (JSON, CSV, text)
- REST API for remote administration

**CLI Examples:**

```bash
# Topic management
streamforge topic list
streamforge topic create events --partitions 3 --replication-factor 3
streamforge topic describe events
streamforge topic delete events

# Message operations
streamforge produce --topic events --key user1 --value "event data"
streamforge consume --topic events --group my-group --max-messages 10

# Start REST API server
streamforge server --port 8080
```

**REST API Examples:**

```bash
# List topics
curl http://localhost:8080/api/v1/topics

# Create topic
curl -X POST http://localhost:8080/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{"name":"events","partitions":3,"replication_factor":3}'

# Produce message
curl -X POST http://localhost:8080/api/v1/produce \
  -H "Content-Type: application/json" \
  -d '{"topic":"events","key":"key1","value":"data"}'
```

**Building:**

```bash
cd admin
go build -o streamforge-admin .

# Run CLI
./streamforge-admin topic list

# Run REST server
./streamforge-admin server --port 8080
```

See [admin/README.md](admin/README.md) for complete CLI and REST API documentation.

### Planned SDKs

- [ ] Python SDK
- [ ] Node.js SDK
- [ ] Java SDK
- [ ] Rust SDK

## Roadmap

- [x] Go SDK (in progress)
- [x] Admin CLI tools (in progress)
- [x] Docker deployment (in progress)
- [x] Performance benchmarks (in progress)
- [ ] Kubernetes operators
- [ ] Advanced compression codecs
- [ ] Schema Registry integration
- [ ] Stream processing framework
- [ ] Additional language SDKs

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues, questions, or suggestions, please open an issue on GitHub.

---

**Status**: Active Development 🚀

