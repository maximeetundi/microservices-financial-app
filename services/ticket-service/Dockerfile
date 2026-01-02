FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY ticket-service/go.mod ticket-service/go.sum* ./
RUN go mod download || true

# Copy source code
COPY ticket-service/ .

# Build - run go mod tidy to update go.sum with new dependencies
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o ticket-service .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/ticket-service .

# Expose port
EXPOSE 8090

# Run
CMD ["./ticket-service"]
