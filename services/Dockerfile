# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy all source code from the 'services' directory
COPY . .

# Tidy and vendor modules
RUN for d in $(cat go.work | grep -E '^\./' | awk '{print $1}'); do (cd "$d" && go mod tidy -e); done
RUN go work vendor

# Build all services
RUN mkdir -p /app/bin
RUN for d in $(cat go.work | grep -E '^\./' | awk '{print $1}'); do (cd "$d" && CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /app/bin/$(basename $d) .); done

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binaries from builder stage
COPY --from=builder /app/bin /app

# Default command (can be overridden in docker-compose.yml)
CMD ["/app/api-gateway"]
