# Stage 1: Build & Tooling
FROM golang:1.24-alpine AS builder

# Set build-time environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on

# Install system dependencies for build and tools
RUN apk add --no-cache \
    git \
    make \
    protobuf-dev \
    curl \
    libc-dev \
    gcc

# Install Go-specific tools (protoc-gen, sqlc, atlas)
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install ariga.io/atlas/cmd/atlas@latest

# Set working directory
WORKDIR /app

# Copy dependency files for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
# -s -w: Omit symbol table and debug information to reduce binary size
RUN go build -ldflags="-s -w" -o goforge ./cmd/main.go

# Stage 2: Production Runtime
FROM alpine:3.21 AS runner

# Install minimal runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create a non-privileged user for security
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/goforge .

# Ensure the binary is executable and owned by appuser
RUN chown appuser:appuser goforge && chmod +x goforge

# Switch to the non-privileged user
USER appuser

# Expose the default application port
EXPOSE 8080

# Standard entrypoint for the microservice
ENTRYPOINT ["./goforge"]

# Default command
CMD ["serve"]
