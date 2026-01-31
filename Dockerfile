# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main ./cmd/api

# Build seeder
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/seeder ./cmd/seeder

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS and wget for health checks
RUN apk --no-cache add ca-certificates wget

# Copy binaries from builder
COPY --from=builder /app/bin/main .
COPY --from=builder /app/bin/seeder .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./main"]
