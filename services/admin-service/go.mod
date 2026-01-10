module github.com/crypto-bank/microservices-financial-app/services/admin-service

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/gin-contrib/cors v1.4.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.4.0
	github.com/lib/pq v1.10.9
	github.com/minio/minio-go/v7 v7.0.66
	github.com/prometheus/client_golang v1.18.0
	github.com/redis/go-redis/v9 v9.3.0
	github.com/segmentio/kafka-go v0.4.47
	golang.org/x/crypto v0.16.0
	github.com/crypto-bank/microservices-financial-app/services/common v0.0.0
)

replace (
    github.com/crypto-bank/microservices-financial-app/services/common => ../common
)
