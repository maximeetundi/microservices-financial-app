# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app/admin-service

COPY ./admin-service .

RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/admin-service/main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/admin-service/main .

EXPOSE 8088

CMD ["./main"]
