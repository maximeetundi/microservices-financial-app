# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/wallet-service

# Copy service source code
COPY ./wallet-service .

# Tidy and vendor modules
RUN go mod tidy -e
RUN go mod vendor

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /app/wallet-service/main .

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