# Build stage
FROM rust:latest as builder

WORKDIR /build

# Install build dependencies
RUN apt-get update && apt-get install -y \
    cmake \
    build-essential \
    pkg-config \
    libssl-dev \
    zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy project files
COPY .. .

# Remove Cargo.lock if it exists (incompatible with older Rust versions)
RUN rm -f /build/broker/Cargo.lock

# Build storage backend (C++)
RUN cd storage && \
    rm -rf build && \
    mkdir build && \
    cd build && \
    cmake .. && \
    cmake --build . --config Release && \
    cd ../..

# Build broker
RUN cd broker && \
    cargo build --release && \
    cd ..

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libssl3 \
    && rm -rf /var/lib/apt/lists/*

# Copy C++ storage library from builder
COPY --from=builder /build/storage/build/libstreamforge_storage_ffi.so /usr/local/lib/
COPY --from=builder /build/storage/build/libstreamforge_storage.a /usr/local/lib/

# Copy compiled broker from builder
COPY --from=builder /build/broker/target/release/broker /app/broker

# Create data directory
RUN mkdir -p /data/streamforge

# Set library path
ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

# Expose broker port
EXPOSE 9092

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/bin/sh", "-c", "nc -z localhost 9092 || exit 1"]

# Run broker
CMD ["/app/broker"]
