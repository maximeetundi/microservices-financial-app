FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY association-service/go.mod association-service/go.sum* ./
RUN go mod download || true

# Copy source code
COPY association-service/ .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o association-service .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/association-service .

# Expose port
EXPOSE 8091

# Run
CMD ["./association-service"]
