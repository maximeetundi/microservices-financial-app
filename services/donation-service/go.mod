module github.com/crypto-bank/microservices-financial-app/services/donation-service

go 1.21

require (
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.3.1
	github.com/minio/minio-go/v7 v7.0.63
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/prometheus/client_golang v1.16.0
	go.mongodb.org/mongo-driver v1.12.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/crypto-bank/microservices-financial-app/services/common v0.0.0-20230101000000-000000000000
)

replace github.com/crypto-bank/microservices-financial-app/services/common => ../common
