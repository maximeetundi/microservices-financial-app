# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go mod
RUN apk add --no-cache git

# Copy common module first (required for local replace directive)
COPY ./common ./common

# Copy go mod files
COPY ./transfer-service/go.mod ./transfer-service/go.sum* ./transfer-service/

# Download dependencies
WORKDIR /app/transfer-service
RUN go mod download || true

# Copy service source code
COPY ./transfer-service .

# Tidy and download dependencies (regenerates go.sum if needed)
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/transfer-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/transfer-service/main .

# Expose port
EXPOSE 8084

# Command to run
CMD ["./main"]