module github.com/crypto-bank/microservices-financial-app/services/enterprise-service

go 1.21

require (
	github.com/crypto-bank/microservices-financial-app/services/common v0.0.0-00010101000000-000000000000
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.1
	go.mongodb.org/mongo-driver v1.13.1
)

replace (
	github.com/crypto-bank/microservices-financial-app/services/common => ../common
)
