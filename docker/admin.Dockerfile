# Build stage
FROM golang:1.21-alpine as builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy admin code
COPY admin .

# Initialize Go module with proper import path
RUN go mod init github.com/iamrahulgupta/streamforge-admin 2>/dev/null || true

# Download dependencies and resolve imports
RUN go mod tidy

# Build admin tool
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o streamforge-admin .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl

# Copy compiled admin from builder
COPY --from=builder /build/streamforge-admin /app/streamforge-admin

# Expose REST API port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Default command to run REST API server
ENTRYPOINT ["/app/streamforge-admin"]
CMD ["server", "--port", "8080", "--brokers", "broker-1:9092,broker-2:9092,broker-3:9092"]
