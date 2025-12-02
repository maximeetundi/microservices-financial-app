# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/exchange-service

# Copy service source code
COPY ./exchange-service .

# Tidy and vendor modules
RUN go mod tidy -e
RUN go mod vendor

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /app/exchange-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/exchange-service/main .

# Expose port
EXPOSE 8085

# Command to run
CMD ["./main"]