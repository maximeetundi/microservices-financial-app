# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for go mod
RUN apk add --no-cache git

# Copy common module first (required for local replace directive)
COPY ./common ./common

# Create service directory
WORKDIR /app/enterprise-service

# Copy go mod files
COPY ./enterprise-service/go.mod ./enterprise-service/go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy service source code
COPY ./enterprise-service .

# Tidy and download dependencies (regenerates go.sum if needed)
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/enterprise-service/main .

EXPOSE 8097
CMD ["./main"]
