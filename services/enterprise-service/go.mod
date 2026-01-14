module github.com/crypto-bank/microservices-financial-app/services/enterprise-service

go 1.21

require (
	github.com/crypto-bank/microservices-financial-app/services/common v0.0.0-00010101000000-000000000000
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	go.mongodb.org/mongo-driver v1.13.1
	golang.org/x/time v0.5.0
	github.com/minio/minio-go/v7 v7.0.66
	github.com/google/uuid v1.6.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
)

replace (
	github.com/crypto-bank/microservices-financial-app/services/common => ../common
)
