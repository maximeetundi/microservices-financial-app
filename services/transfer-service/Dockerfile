# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/transfer-service

# Copy service source code
COPY ./transfer-service .

# Download dependencies
RUN go mod download

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