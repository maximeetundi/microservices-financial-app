package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/config"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/routes"
	"github.com/crypto-bank/microservices-financial-app/services/api-gateway/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize services
	serviceManager := services.NewServiceManager(cfg)

	// Setup Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS configuration plus permissive pour développement
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001", 
			"https://*.cryptobank.app",
			"https://cryptobank.com",
			"http://127.0.0.1:3000",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization", 
			"X-Requested-With", "X-API-Key", "X-User-ID",
		},
		ExposeHeaders: []string{
			"Content-Length", "Authorization", "X-Total-Count",
		},
		AllowCredentials: true,
		MaxAge: 86400, // 24 hours
	}))

	// Security middleware
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RateLimiter())

	// Health check global
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":      "ok",
			"service":     "api-gateway",
			"version":     "2.0.0",
			"gateway_url": "http://localhost:8080",
			"services": map[string]string{
				"auth":         "/gateway/auth-service",
				"wallets":      "/gateway/wallet-service", 
				"exchange":     "/gateway/exchange-service",
				"cards":        "/gateway/card-service",
				"transfers":    "/gateway/transfer-service",
				"users":        "/gateway/user-service",
				"notifications": "/gateway/notification-service",
			},
		})
	})

	// Routes principales avec préfixe /gateway
	gateway := router.Group("/gateway")
	{
		// Service discovery et documentation
		routes.SetupServiceDiscovery(gateway)
		
		// Routes des services avec préfixes
		routes.SetupServiceRoutes(gateway, serviceManager)
		
		// Routes d'administration
		admin := gateway.Group("/admin")
		admin.Use(middleware.JWTAuth(cfg.JWTSecret))
		admin.Use(middleware.AdminOnly())
		{
			routes.SetupAdminRoutes(admin, serviceManager)
		}
	}

	// Routes de compatibilité (anciennes routes sans préfixe)
	// Pour la transition, on garde les anciennes routes
	legacyApi := router.Group("/api/v1")
	{
		// Auth routes (no auth required)
		auth := legacyApi.Group("/auth")
		routes.SetupAuthRoutes(auth, serviceManager)

		// Protected routes avec middleware
		protected := legacyApi.Group("/")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			routes.SetupUserRoutes(protected.Group("/users"), serviceManager)
			routes.SetupWalletRoutes(protected.Group("/wallets"), serviceManager)
			routes.SetupTransferRoutes(protected.Group("/transfers"), serviceManager)
			routes.SetupExchangeRoutes(protected.Group("/exchange"), serviceManager)
			routes.SetupCardRoutes(protected.Group("/cards"), serviceManager)
			routes.SetupNotificationRoutes(protected.Group("/notifications"), serviceManager)
		}

		// Admin routes legacy
		adminLegacy := legacyApi.Group("/admin")
		adminLegacy.Use(middleware.JWTAuth(cfg.JWTSecret))
		adminLegacy.Use(middleware.AdminOnly())
		{
			routes.SetupAdminRoutes(adminLegacy, serviceManager)
		}
	}

	// Page d'accueil avec documentation interactive
	router.GET("/", func(c *gin.Context) {
		html := `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CryptoBank API Gateway</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #2c3e50; margin-bottom: 10px; }
        .subtitle { color: #7f8c8d; margin-bottom: 30px; }
        .services { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin: 30px 0; }
        .service { background: #ecf0f1; padding: 20px; border-radius: 8px; border-left: 4px solid #3498db; }
        .service h3 { margin: 0 0 10px 0; color: #2c3e50; }
        .service p { margin: 0 0 15px 0; color: #7f8c8d; font-size: 14px; }
        .endpoints { list-style: none; padding: 0; }
        .endpoints li { background: #fff; padding: 8px 12px; margin: 5px 0; border-radius: 4px; font-family: monospace; font-size: 12px; }
        .url-example { background: #2c3e50; color: white; padding: 15px; border-radius: 8px; margin: 20px 0; font-family: monospace; }
        .footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #ecf0f1; text-align: center; color: #7f8c8d; }
        .status-ok { display: inline-block; background: #27ae60; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🏦 CryptoBank API Gateway</h1>
        <p class="subtitle">Passerelle API unifiée pour tous les services bancaires crypto</p>
        
        <div class="url-example">
            🌍 Gateway Base URL: <strong>http://localhost:8080/gateway/{service-name}/</strong>
        </div>

        <div class="services">
            <div class="service">
                <h3>🔐 Auth Service <span class="status-ok">ACTIVE</span></h3>
                <p>Authentification, inscription, 2FA, gestion des sessions</p>
                <ul class="endpoints">
                    <li>POST /gateway/auth-service/login</li>
                    <li>POST /gateway/auth-service/register</li>
                    <li>POST /gateway/auth-service/enable-2fa</li>
                </ul>
            </div>

            <div class="service">
                <h3>💰 Wallet Service <span class="status-ok">ACTIVE</span></h3>
                <p>Portefeuilles crypto et fiat, transactions</p>
                <ul class="endpoints">
                    <li>GET /gateway/wallet-service/wallets</li>
                    <li>POST /gateway/wallet-service/crypto/generate</li>
                    <li>POST /gateway/wallet-service/crypto/:id/send</li>
                </ul>
            </div>

            <div class="service">
                <h3>🔄 Exchange Service <span class="status-ok">ACTIVE</span></h3>
                <p>Trading crypto, conversion fiat, taux de change</p>
                <ul class="endpoints">
                    <li>GET /gateway/exchange-service/rates</li>
                    <li>POST /gateway/exchange-service/trading/buy</li>
                    <li>POST /gateway/exchange-service/fiat/convert</li>
                </ul>
            </div>

            <div class="service">
                <h3>💳 Card Service <span class="status-ok">ACTIVE</span></h3>
                <p>Cartes prépayées, gift cards, cartes virtuelles/physiques</p>
                <ul class="endpoints">
                    <li>POST /gateway/card-service/virtual/</li>
                    <li>POST /gateway/card-service/physical/</li>
                    <li>POST /gateway/card-service/gift/</li>
                </ul>
            </div>

            <div class="service">
                <h3>💸 Transfer Service <span class="status-ok">ACTIVE</span></h3>
                <p>Transferts d'argent, mobile money, international</p>
                <ul class="endpoints">
                    <li>POST /gateway/transfer-service/transfers</li>
                    <li>POST /gateway/transfer-service/mobile/send</li>
                    <li>POST /gateway/transfer-service/international</li>
                </ul>
            </div>

            <div class="service">
                <h3>👤 User Service <span class="status-ok">ACTIVE</span></h3>
                <p>Profils utilisateur, KYC, paramètres</p>
                <ul class="endpoints">
                    <li>GET /gateway/user-service/profile</li>
                    <li>POST /gateway/user-service/kyc/upload</li>
                    <li>PUT /gateway/user-service/settings</li>
                </ul>
            </div>
        </div>

        <div class="service">
            <h3>🚀 Exemples d'utilisation</h3>
            <div style="background: #f8f9fa; padding: 15px; border-radius: 8px; margin: 10px 0;">
                <strong>1. Acheter du Bitcoin :</strong><br>
                <code>POST /gateway/exchange-service/trading/buy</code><br><br>
                
                <strong>2. Convertir USD → EUR :</strong><br>
                <code>POST /gateway/exchange-service/fiat/execute</code><br><br>
                
                <strong>3. Créer une carte virtuelle :</strong><br>
                <code>POST /gateway/card-service/virtual/</code><br><br>
                
                <strong>4. Transfert mobile money :</strong><br>
                <code>POST /gateway/transfer-service/mobile/send</code>
            </div>
        </div>

        <div class="footer">
            <p>📚 <a href="/gateway/docs">Documentation complète</a> | 
               🔍 <a href="/gateway/services">Liste des services</a> | 
               💻 <a href="http://localhost:3000">Interface Web</a></p>
            <p>CryptoBank API Gateway v2.0.0 - Microservices Architecture</p>
        </div>
    </div>

    <script>
        // Auto-refresh status every 30 seconds
        setInterval(() => {
            fetch('/health')
                .then(r => r.json())
                .then(data => console.log('Gateway Health:', data.status))
                .catch(e => console.error('Gateway offline'));
        }, 30000);
    </script>
</body>
</html>`
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, html)
	})

	// Page de test simple
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "CryptoBank API Gateway Test",
			"routes": map[string]string{
				"service_discovery": "/gateway/services",
				"documentation":     "/gateway/docs", 
				"health":           "/health",
				"legacy_api":       "/api/v1/*",
				"new_gateway":      "/gateway/{service}/*",
			},
			"examples": map[string]map[string]interface{}{
				"fiat_conversion": {
					"url": "/gateway/exchange-service/fiat/convert?from=USD&to=EUR&amount=1000",
					"description": "Convertir 1000 USD en EUR",
				},
				"crypto_purchase": {
					"url": "/gateway/exchange-service/trading/buy",
					"method": "POST",
					"description": "Acheter des cryptomonnaies",
				},
				"card_creation": {
					"url": "/gateway/card-service/virtual/",
					"method": "POST", 
					"description": "Créer une carte virtuelle",
				},
			},
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 CryptoBank API Gateway v2.0 starting on port %s", port)
	log.Printf("📍 Gateway URL: http://localhost:%s/gateway/", port)
	log.Printf("📚 Documentation: http://localhost:%s/gateway/docs", port)
	log.Printf("🌐 Web Interface: http://localhost:3000")
	log.Printf("📊 Health Check: http://localhost:%s/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}