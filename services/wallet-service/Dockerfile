# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/wallet-service

# Install git for go mod
RUN apk add --no-cache git

# Copy go mod files
COPY ./wallet-service/go.mod ./wallet-service/go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy service source code
COPY ./wallet-service .

# Tidy and download dependencies (regenerates go.sum if needed)
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/wallet-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/wallet-service/main .

# Expose port
EXPOSE 8082

# Command to run
CMD ["./main"]