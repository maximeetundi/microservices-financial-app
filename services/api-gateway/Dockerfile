# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/api-gateway

# Install git for go mod
RUN apk add --no-cache git

# Copy go mod files
COPY ./api-gateway/go.mod ./api-gateway/go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy service source code
COPY ./api-gateway .

# Tidy and download dependencies (regenerates go.sum if needed)
RUN go mod tidy && go mod download

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