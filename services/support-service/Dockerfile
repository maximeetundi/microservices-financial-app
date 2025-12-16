# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy go mod files
COPY ./support-service/go.mod ./support-service/go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy service source code
COPY ./support-service .

# Tidy and download dependencies
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/support-service/main .

# Final stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/support-service/main .

EXPOSE 8089

CMD ["./main"]
