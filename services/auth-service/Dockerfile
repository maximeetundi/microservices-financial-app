# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/auth-service

# Install git for go mod
RUN apk add --no-cache git

# Copy go mod files
COPY ./auth-service/go.mod ./auth-service/go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy service source code
COPY ./auth-service .

# Tidy and download dependencies (regenerates go.sum if needed)
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