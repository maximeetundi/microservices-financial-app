# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/api-gateway

# Copy service source code
COPY ./api-gateway .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-gateway/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/api-gateway/main .

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]