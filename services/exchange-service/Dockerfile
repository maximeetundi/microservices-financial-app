# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/exchange-service

# Install git for go mod
RUN apk add --no-cache git

# Copy go mod files
COPY ./exchange-service/go.mod ./exchange-service/go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy service source code
COPY ./exchange-service .

# Tidy and download dependencies (regenerates go.sum if needed)
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/exchange-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/exchange-service/main .

# Expose port
EXPOSE 8084

# Command to run
CMD ["./main"]