# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/card-service

# Install git for go mod
RUN apk add --no-cache git

# Copy go mod files
COPY ./card-service/go.mod ./card-service/go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy service source code
COPY ./card-service .

# Tidy and download dependencies (regenerates go.sum if needed)
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/card-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/card-service/main .

# Expose port
EXPOSE 8085

# Command to run
CMD ["./main"]