# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go mod
RUN apk add --no-cache git

# Copy common module first (required for local replace directive)
COPY ./common ./common

# Copy go mod file
COPY ./wallet-service/go.mod ./wallet-service/

# Download dependencies
WORKDIR /app/wallet-service
RUN go mod download || true

# Copy service source code
COPY ./wallet-service .

# Pin dependencies to stable versions compatible with Go 1.23
RUN go get github.com/ethereum/go-ethereum@v1.13.15

# Tidy dependencies after copying source
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/wallet-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/wallet-service/main .

# Expose port
EXPOSE 8083

# Command to run
CMD ["./main"]