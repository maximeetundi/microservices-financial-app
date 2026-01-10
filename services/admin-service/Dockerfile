# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for go mod
RUN apk add --no-cache git

# Copy common module first (required for local replace directive)
COPY ./common ./common

# Copy go mod files
COPY ./admin-service/go.mod ./admin-service/go.sum* ./admin-service/

# Download dependencies
WORKDIR /app/admin-service
RUN go mod download || true

# Copy service source code
COPY ./admin-service .

# Tidy and download dependencies (regenerates go.sum if needed)
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/admin-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/admin-service/main .

# Expose port
EXPOSE 8088

# Command to run
CMD ["./main"]
