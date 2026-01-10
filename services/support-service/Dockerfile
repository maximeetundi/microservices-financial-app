# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go mod
RUN apk add --no-cache git

# Copy common module first (required for local replace directive)
COPY ./common ./common

# Copy go mod files
COPY ./support-service/go.mod ./support-service/go.sum* ./support-service/

# Download dependencies
WORKDIR /app/support-service
RUN go mod download || true

# Copy source code
COPY ./support-service .

# Build
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o support-service .

# Runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

# Copy binary
COPY --from=builder /app/support-service/support-service .

# Expose port
EXPOSE 8089

# Run
CMD ["./support-service"]
