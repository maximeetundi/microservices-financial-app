# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/auth-service

# Copy service source code
COPY ./auth-service .

# Download dependencies
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/auth-service/main .

# Expose port
EXPOSE 8081

# Command to run
CMD ["./main"]