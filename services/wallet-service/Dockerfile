# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy all source code from the 'services' directory
COPY . .

# Run go work sync and vendor from the workspace root
RUN for d in $(cat go.work | grep -E '^\\s*\\./' | awk '{print $1}'); do (cd "$d" && /usr/local/go/bin/go mod tidy -e); done
RUN /usr/local/go/bin/go work vendor

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o main ./wallet-service

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8083

# Command to run
CMD ["./main"]