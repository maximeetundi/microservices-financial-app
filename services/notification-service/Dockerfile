# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/notification-service

# Copy service source code
COPY ./notification-service .

# Download dependencies (go mod tidy generates go.sum if missing)
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notification-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/notification-service/main .

# Expose port
EXPOSE 8087

# Command to run
CMD ["./main"]
