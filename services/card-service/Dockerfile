# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/card-service

# Copy service source code
COPY ./card-service .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/card-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/card-service/main .

# Expose port
EXPOSE 8086

# Command to run
CMD ["./main"]