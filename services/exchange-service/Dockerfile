# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go mod
RUN apk add --no-cache git

# Copy common module first (required for local replace directive)
COPY ./common ./common

# Copy go mod files
# Download dependencies
WORKDIR /app/exchange-service
RUN go mod download || true

# Copy go mod file
COPY ./exchange-service/go.mod ./

# Copy service source code
COPY ./exchange-service .

# Tidy dependencies after copying source
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/exchange-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/exchange-service/main .

# Expose port
EXPOSE 8085

# Command to run
CMD ["./main"]