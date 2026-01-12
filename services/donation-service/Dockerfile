FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Download all dependencies. Run tidy to ensure go.sum exists
RUN go mod tidy
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o donation-service .

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/donation-service .

# Expose port 8096 to the outside world
EXPOSE 8096

# Command to run the executable
CMD ["./donation-service"]
