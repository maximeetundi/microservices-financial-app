module github.com/crypto-bank/microservices-financial-app/services/shop-service

go 1.21

require (
	github.com/crypto-bank/microservices-financial-app/services/common v0.0.0
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/minio/minio-go/v7 v7.0.63
	github.com/prometheus/client_golang v1.17.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	go.mongodb.org/mongo-driver v1.13.1
)

replace github.com/crypto-bank/microservices-financial-app/services/common => ../common
