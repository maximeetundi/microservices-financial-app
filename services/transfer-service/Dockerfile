# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/transfer-service

# Install git for go mod
RUN apk add --no-cache git

# Copy go mod files
COPY ./transfer-service/go.mod ./transfer-service/go.sum* ./

# Download dependencies
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
EXPOSE 8083

# Command to run
CMD ["./main"]