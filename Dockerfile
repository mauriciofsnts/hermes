# ========================================
# Stage 1: Build
# ========================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
# -ldflags: strip debug info and reduce binary size
# CGO_ENABLED=0: static binary without C dependencies
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -gcflags=all=-l \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty)" \
    -o hermes \
    ./cmd/hermes

# ========================================
# Stage 2: Runtime
# ========================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && update-ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 hermes && \
    adduser -D -u 1000 -G hermes hermes

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=hermes:hermes /app/hermes /app/hermes

# Copy configuration template (can be overridden via volume)
COPY --chown=hermes:hermes config_example.yaml /app/config.yaml

# Copy templates directory
COPY --chown=hermes:hermes templates /app/templates

# Switch to non-root user
USER hermes

# Expose HTTP port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set entrypoint
ENTRYPOINT ["/app/hermes"]
