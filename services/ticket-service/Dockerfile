FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download || true

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o ticket-service .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/ticket-service .

# Expose port
EXPOSE 8090

# Run
CMD ["./ticket-service"]
