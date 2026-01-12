FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy common module first (required for local replace directive)
COPY ./common ./common

# Copy go mod file
COPY ./donation-service/go.mod ./donation-service/

# Download dependencies
WORKDIR /app/donation-service
RUN go mod download || true

# Copy the source from the current directory to the Working Directory inside the container
COPY ./donation-service .

# Tidy dependencies after copying source
RUN go mod tidy

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o donation-service .

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/donation-service/donation-service .

# Expose port 8096 to the outside world
EXPOSE 8096

# Command to run the executable
CMD ["./donation-service"]
