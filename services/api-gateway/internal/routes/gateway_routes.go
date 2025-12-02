package routes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
)

// Configuration des routes par service avec préfixe
func SetupServiceRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	// Routes préfixées par service pour éviter les conflits d'URL
	
	// Auth Service Routes - /gateway/auth-service/*
	authGroup := router.Group("/auth-service")
	{
		authGroup.POST("/register", proxyToService(serviceManager, "auth", "/api/v1/register"))
		authGroup.POST("/login", proxyToService(serviceManager, "auth", "/api/v1/login"))
		authGroup.POST("/refresh", proxyToService(serviceManager, "auth", "/api/v1/refresh"))
		authGroup.POST("/logout", proxyToService(serviceManager, "auth", "/api/v1/logout"))
		authGroup.POST("/forgot-password", proxyToService(serviceManager, "auth", "/api/v1/forgot-password"))
		authGroup.POST("/reset-password", proxyToService(serviceManager, "auth", "/api/v1/reset-password"))
		authGroup.POST("/verify-email", proxyToService(serviceManager, "auth", "/api/v1/verify-email"))
		authGroup.POST("/verify-phone", proxyToService(serviceManager, "auth", "/api/v1/verify-phone"))
		authGroup.POST("/enable-2fa", proxyToService(serviceManager, "auth", "/api/v1/enable-2fa"))
		authGroup.POST("/verify-2fa", proxyToService(serviceManager, "auth", "/api/v1/verify-2fa"))
		authGroup.GET("/sessions", proxyToService(serviceManager, "auth", "/api/v1/sessions"))
		authGroup.DELETE("/sessions/:session_id", proxyToService(serviceManager, "auth", "/api/v1/sessions/:session_id"))
	}

	// Wallet Service Routes - /gateway/wallet-service/*
	walletGroup := router.Group("/wallet-service")
	{
		walletGroup.GET("/wallets", proxyToService(serviceManager, "wallet", "/api/v1/wallets"))
		walletGroup.POST("/wallets", proxyToService(serviceManager, "wallet", "/api/v1/wallets"))
		walletGroup.GET("/wallets/:wallet_id", proxyToService(serviceManager, "wallet", "/api/v1/wallets/:wallet_id"))
		walletGroup.PUT("/wallets/:wallet_id", proxyToService(serviceManager, "wallet", "/api/v1/wallets/:wallet_id"))
		walletGroup.GET("/wallets/:wallet_id/balance", proxyToService(serviceManager, "wallet", "/api/v1/wallets/:wallet_id/balance"))
		walletGroup.GET("/wallets/:wallet_id/transactions", proxyToService(serviceManager, "wallet", "/api/v1/wallets/:wallet_id/transactions"))
		walletGroup.POST("/wallets/:wallet_id/freeze", proxyToService(serviceManager, "wallet", "/api/v1/wallets/:wallet_id/freeze"))
		walletGroup.POST("/wallets/:wallet_id/unfreeze", proxyToService(serviceManager, "wallet", "/api/v1/wallets/:wallet_id/unfreeze"))
		
		// Crypto wallet specific
		cryptoGroup := walletGroup.Group("/crypto")
		{
			cryptoGroup.POST("/generate", proxyToService(serviceManager, "wallet", "/api/v1/crypto/generate"))
			cryptoGroup.GET("/:wallet_id/address", proxyToService(serviceManager, "wallet", "/api/v1/crypto/:wallet_id/address"))
			cryptoGroup.POST("/:wallet_id/send", proxyToService(serviceManager, "wallet", "/api/v1/crypto/:wallet_id/send"))
			cryptoGroup.GET("/:wallet_id/pending", proxyToService(serviceManager, "wallet", "/api/v1/crypto/:wallet_id/pending"))
			cryptoGroup.POST("/:wallet_id/estimate-fee", proxyToService(serviceManager, "wallet", "/api/v1/crypto/:wallet_id/estimate-fee"))
		}
	}

	// Transfer Service Routes - /gateway/transfer-service/*
	transferGroup := router.Group("/transfer-service")
	{
		transferGroup.POST("/transfers", proxyToService(serviceManager, "transfer", "/api/v1/transfers"))
		transferGroup.GET("/transfers", proxyToService(serviceManager, "transfer", "/api/v1/transfers"))
		transferGroup.GET("/transfers/:transfer_id", proxyToService(serviceManager, "transfer", "/api/v1/transfers/:transfer_id"))
		transferGroup.POST("/transfers/:transfer_id/cancel", proxyToService(serviceManager, "transfer", "/api/v1/transfers/:transfer_id/cancel"))
		transferGroup.POST("/international", proxyToService(serviceManager, "transfer", "/api/v1/international"))
		
		// Mobile money
		mobileGroup := transferGroup.Group("/mobile")
		{
			mobileGroup.POST("/send", proxyToService(serviceManager, "transfer", "/api/v1/mobile/send"))
			mobileGroup.POST("/receive", proxyToService(serviceManager, "transfer", "/api/v1/mobile/receive"))
			mobileGroup.GET("/providers", proxyToService(serviceManager, "transfer", "/api/v1/mobile/providers"))
		}
		
		// Bulk transfers
		bulkGroup := transferGroup.Group("/bulk")
		{
			bulkGroup.POST("/", proxyToService(serviceManager, "transfer", "/api/v1/bulk"))
			bulkGroup.GET("/:batch_id", proxyToService(serviceManager, "transfer", "/api/v1/bulk/:batch_id"))
			bulkGroup.POST("/:batch_id/approve", proxyToService(serviceManager, "transfer", "/api/v1/bulk/:batch_id/approve"))
		}
	}

	// Exchange Service Routes - /gateway/exchange-service/*
	exchangeGroup := router.Group("/exchange-service")
	{
		// Taux de change publics (pas d'auth requis)
		exchangeGroup.GET("/rates", proxyToService(serviceManager, "exchange", "/api/v1/rates"))
		exchangeGroup.GET("/rates/:from/:to", proxyToService(serviceManager, "exchange", "/api/v1/rates/:from/:to"))
		exchangeGroup.GET("/markets", proxyToService(serviceManager, "exchange", "/api/v1/markets"))
		exchangeGroup.GET("/orderbook/:pair", proxyToService(serviceManager, "exchange", "/api/v1/orderbook/:pair"))
		exchangeGroup.GET("/trades/:pair", proxyToService(serviceManager, "exchange", "/api/v1/trades/:pair"))

		// Exchange operations (auth requis)
		exchangeGroup.POST("/quote", proxyToService(serviceManager, "exchange", "/api/v1/quote"))
		exchangeGroup.POST("/execute", proxyToService(serviceManager, "exchange", "/api/v1/execute"))
		exchangeGroup.GET("/history", proxyToService(serviceManager, "exchange", "/api/v1/history"))
		exchangeGroup.GET("/exchange/:exchange_id", proxyToService(serviceManager, "exchange", "/api/v1/exchange/:exchange_id"))

		// Trading operations
		tradingGroup := exchangeGroup.Group("/trading")
		{
			tradingGroup.POST("/buy", proxyToService(serviceManager, "exchange", "/api/v1/trading/buy"))
			tradingGroup.POST("/sell", proxyToService(serviceManager, "exchange", "/api/v1/trading/sell"))
			tradingGroup.POST("/limit-order", proxyToService(serviceManager, "exchange", "/api/v1/trading/limit-order"))
			tradingGroup.POST("/stop-loss", proxyToService(serviceManager, "exchange", "/api/v1/trading/stop-loss"))
			tradingGroup.GET("/orders", proxyToService(serviceManager, "exchange", "/api/v1/trading/orders"))
			tradingGroup.GET("/orders/active", proxyToService(serviceManager, "exchange", "/api/v1/trading/orders/active"))
			tradingGroup.DELETE("/orders/:order_id", proxyToService(serviceManager, "exchange", "/api/v1/trading/orders/:order_id"))
			tradingGroup.GET("/portfolio", proxyToService(serviceManager, "exchange", "/api/v1/trading/portfolio"))
			tradingGroup.GET("/performance", proxyToService(serviceManager, "exchange", "/api/v1/trading/performance"))
		}

		// Conversion fiat-fiat
		fiatGroup := exchangeGroup.Group("/fiat")
		{
			fiatGroup.GET("/rates", proxyToService(serviceManager, "exchange", "/api/v1/fiat/rates"))
			fiatGroup.GET("/rates/:from/:to", proxyToService(serviceManager, "exchange", "/api/v1/fiat/rates/:from/:to"))
			fiatGroup.GET("/convert", proxyToService(serviceManager, "exchange", "/api/v1/fiat/convert"))
			fiatGroup.POST("/quote", proxyToService(serviceManager, "exchange", "/api/v1/fiat/quote"))
			fiatGroup.POST("/execute", proxyToService(serviceManager, "exchange", "/api/v1/fiat/execute"))
			fiatGroup.GET("/currencies", proxyToService(serviceManager, "exchange", "/api/v1/fiat/currencies"))
			fiatGroup.GET("/compare-fees", proxyToService(serviceManager, "exchange", "/api/v1/fiat/compare-fees"))
		}

		// P2P Trading
		p2pGroup := exchangeGroup.Group("/p2p")
		{
			p2pGroup.POST("/offers", proxyToService(serviceManager, "exchange", "/api/v1/p2p/offers"))
			p2pGroup.GET("/offers", proxyToService(serviceManager, "exchange", "/api/v1/p2p/offers"))
			p2pGroup.POST("/offers/:offer_id/accept", proxyToService(serviceManager, "exchange", "/api/v1/p2p/offers/:offer_id/accept"))
			p2pGroup.GET("/trades", proxyToService(serviceManager, "exchange", "/api/v1/p2p/trades"))
		}
	}

	// Card Service Routes - /gateway/card-service/*
	cardGroup := router.Group("/card-service")
	{
		cardGroup.GET("/cards", proxyToService(serviceManager, "card", "/api/v1/cards"))
		cardGroup.POST("/cards", proxyToService(serviceManager, "card", "/api/v1/cards"))
		cardGroup.GET("/cards/:card_id", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id"))
		cardGroup.PUT("/cards/:card_id", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id"))
		cardGroup.DELETE("/cards/:card_id", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id"))

		// Card operations
		cardGroup.POST("/cards/:card_id/activate", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/activate"))
		cardGroup.POST("/cards/:card_id/deactivate", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/deactivate"))
		cardGroup.POST("/cards/:card_id/freeze", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/freeze"))
		cardGroup.POST("/cards/:card_id/unfreeze", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/unfreeze"))
		cardGroup.POST("/cards/:card_id/block", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/block"))

		// Card loading
		cardGroup.POST("/cards/:card_id/load", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/load"))
		cardGroup.POST("/cards/:card_id/auto-load", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/auto-load"))
		cardGroup.DELETE("/cards/:card_id/auto-load", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/auto-load"))

		// Card details
		cardGroup.GET("/cards/:card_id/limits", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/limits"))
		cardGroup.PUT("/cards/:card_id/limits", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/limits"))
		cardGroup.GET("/cards/:card_id/transactions", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/transactions"))
		cardGroup.GET("/cards/:card_id/balance", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/balance"))
		cardGroup.GET("/cards/:card_id/details", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/details"))
		cardGroup.POST("/cards/:card_id/pin", proxyToService(serviceManager, "card", "/api/v1/cards/:card_id/pin"))

		// Virtual cards
		virtualGroup := cardGroup.Group("/virtual")
		{
			virtualGroup.POST("/", proxyToService(serviceManager, "card", "/api/v1/cards/virtual"))
			virtualGroup.POST("/:card_id/regenerate", proxyToService(serviceManager, "card", "/api/v1/cards/virtual/:card_id/regenerate"))
			virtualGroup.GET("/:card_id/qr", proxyToService(serviceManager, "card", "/api/v1/cards/virtual/:card_id/qr"))
		}

		// Physical cards
		physicalGroup := cardGroup.Group("/physical")
		{
			physicalGroup.POST("/", proxyToService(serviceManager, "card", "/api/v1/cards/physical"))
			physicalGroup.GET("/:card_id/shipping", proxyToService(serviceManager, "card", "/api/v1/cards/physical/:card_id/shipping"))
			physicalGroup.POST("/:card_id/activate", proxyToService(serviceManager, "card", "/api/v1/cards/physical/:card_id/activate-physical"))
			physicalGroup.POST("/:card_id/report-lost", proxyToService(serviceManager, "card", "/api/v1/cards/physical/:card_id/report-lost"))
			physicalGroup.POST("/:card_id/replacement", proxyToService(serviceManager, "card", "/api/v1/cards/physical/:card_id/replacement"))
		}

		// Gift cards
		giftGroup := cardGroup.Group("/gift")
		{
			giftGroup.POST("/", proxyToService(serviceManager, "card", "/api/v1/cards/gift"))
			giftGroup.POST("/send", proxyToService(serviceManager, "card", "/api/v1/cards/gift/send"))
			giftGroup.POST("/redeem", proxyToService(serviceManager, "card", "/api/v1/cards/gift/redeem"))
			giftGroup.GET("/", proxyToService(serviceManager, "card", "/api/v1/cards/gift"))
		}

		// Public card info
		cardGroup.GET("/supported-currencies", proxyToService(serviceManager, "card", "/api/v1/public/supported-currencies"))
		cardGroup.GET("/fees", proxyToService(serviceManager, "card", "/api/v1/public/fees"))
	}

	// User Service Routes - /gateway/user-service/*
	userGroup := router.Group("/user-service")
	{
		userGroup.GET("/profile", proxyToService(serviceManager, "user", "/api/v1/users/profile"))
		userGroup.PUT("/profile", proxyToService(serviceManager, "user", "/api/v1/users/profile"))
		userGroup.GET("/kyc", proxyToService(serviceManager, "user", "/api/v1/users/kyc"))
		userGroup.POST("/kyc/upload", proxyToService(serviceManager, "user", "/api/v1/users/kyc/upload"))
		userGroup.GET("/settings", proxyToService(serviceManager, "user", "/api/v1/users/settings"))
		userGroup.PUT("/settings", proxyToService(serviceManager, "user", "/api/v1/users/settings"))
	}

	// Notification Service Routes - /gateway/notification-service/*
	notificationGroup := router.Group("/notification-service")
	{
		notificationGroup.GET("/notifications", proxyToService(serviceManager, "notification", "/api/v1/notifications"))
		notificationGroup.PUT("/notifications/:notification_id/read", proxyToService(serviceManager, "notification", "/api/v1/notifications/:notification_id/read"))
		notificationGroup.DELETE("/notifications/:notification_id", proxyToService(serviceManager, "notification", "/api/v1/notifications/:notification_id"))
		notificationGroup.POST("/settings", proxyToService(serviceManager, "notification", "/api/v1/notifications/settings"))
	}
}

// Fonction de proxy vers les services
func proxyToService(serviceManager *services.ServiceManager, serviceName, targetPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Construire la nouvelle URL en remplaçant les paramètres
		finalPath := replacePlaceholders(targetPath, c)
		
		// Ajouter les query parameters
		if c.Request.URL.RawQuery != "" {
			finalPath += "?" + c.Request.URL.RawQuery
		}

		// Extraire les headers nécessaires
		headers := make(map[string]string)
		
		// Transférer les headers d'authentification
		if auth := c.GetHeader("Authorization"); auth != "" {
			headers["Authorization"] = auth
		}
		
		// Transférer d'autres headers importants
		if contentType := c.GetHeader("Content-Type"); contentType != "" {
			headers["Content-Type"] = contentType
		}
		
		if userAgent := c.GetHeader("User-Agent"); userAgent != "" {
			headers["User-Agent"] = userAgent
		}

		// Lire le body de la requête si présent
		var requestBody interface{}
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.Body != nil {
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err == nil && len(bodyBytes) > 0 {
					// Remettre le body dans la requête pour les prochains middlewares
					c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					
					// Parser le JSON si possible
					var jsonBody interface{}
					if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
						requestBody = jsonBody
					}
				}
			}
		}

		// Appeler le service
		response, err := serviceManager.CallService(
			c.Request.Context(),
			serviceName,
			c.Request.Method,
			finalPath,
			requestBody,
			headers,
		)

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   "Service unavailable",
				"service": serviceName,
				"message": err.Error(),
			})
			return
		}

		// Transférer les headers de réponse
		for key, values := range response.Headers {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Retourner la réponse
		c.Data(response.StatusCode, "application/json", response.Body)
	}
}

// Remplacer les placeholders dans le path (:param)
func replacePlaceholders(path string, c *gin.Context) string {
	result := path
	
	// Remplacer les paramètres de route
	for _, param := range c.Params {
		placeholder := ":" + param.Key
		result = strings.ReplaceAll(result, placeholder, param.Value)
	}
	
	return result
}

// Route pour obtenir la liste des services disponibles
func SetupServiceDiscovery(router *gin.RouterGroup) {
	router.GET("/services", func(c *gin.Context) {
		services := map[string]interface{}{
			"auth-service": map[string]interface{}{
				"name":        "Authentication Service",
				"description": "User authentication, registration, 2FA",
				"base_path":   "/gateway/auth-service",
				"status":      "active",
				"endpoints": []string{
					"POST /register", "POST /login", "POST /logout",
					"POST /enable-2fa", "POST /verify-2fa",
				},
			},
			"wallet-service": map[string]interface{}{
				"name":        "Wallet Service",
				"description": "Crypto and fiat wallet management",
				"base_path":   "/gateway/wallet-service",
				"status":      "active",
				"endpoints": []string{
					"GET /wallets", "POST /wallets",
					"POST /crypto/generate", "POST /crypto/:id/send",
				},
			},
			"exchange-service": map[string]interface{}{
				"name":        "Exchange Service", 
				"description": "Crypto trading, fiat conversion, exchange rates",
				"base_path":   "/gateway/exchange-service",
				"status":      "active",
				"endpoints": []string{
					"GET /rates", "POST /trading/buy", "POST /trading/sell",
					"POST /fiat/quote", "POST /fiat/execute",
				},
			},
			"card-service": map[string]interface{}{
				"name":        "Card Service",
				"description": "Prepaid cards, gift cards, virtual/physical cards",
				"base_path":   "/gateway/card-service", 
				"status":      "active",
				"endpoints": []string{
					"GET /cards", "POST /cards", "POST /virtual/",
					"POST /physical/", "POST /gift/",
				},
			},
			"transfer-service": map[string]interface{}{
				"name":        "Transfer Service",
				"description": "Money transfers, mobile money, international transfers",
				"base_path":   "/gateway/transfer-service",
				"status":      "active",
				"endpoints": []string{
					"POST /transfers", "POST /international",
					"POST /mobile/send", "POST /bulk/",
				},
			},
		}

		c.JSON(http.StatusOK, gin.H{
			"gateway_version": "1.0.0",
			"services":        services,
			"total_services":  len(services),
			"documentation":   "/gateway/docs",
			"health_check":    "/gateway/health",
		})
	})

	// Documentation des endpoints
	router.GET("/docs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "CryptoBank API Gateway Documentation",
			"version": "1.0.0",
			"services": map[string]string{
				"auth":         "/gateway/auth-service/*",
				"wallets":      "/gateway/wallet-service/*",
				"exchange":     "/gateway/exchange-service/*",
				"cards":        "/gateway/card-service/*",
				"transfers":    "/gateway/transfer-service/*",
				"users":        "/gateway/user-service/*",
				"notifications": "/gateway/notification-service/*",
			},
			"examples": map[string]string{
				"login":          "POST /gateway/auth-service/login",
				"get_wallets":    "GET /gateway/wallet-service/wallets",
				"buy_crypto":     "POST /gateway/exchange-service/trading/buy",
				"convert_fiat":   "POST /gateway/exchange-service/fiat/execute",
				"create_card":    "POST /gateway/card-service/cards",
				"transfer_money": "POST /gateway/transfer-service/transfers",
			},
		})
	})
}