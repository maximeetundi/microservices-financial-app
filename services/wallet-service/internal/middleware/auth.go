package middleware

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	KYCLevel int    `json:"kyc_level"`
	jwt.RegisteredClaims
}

var (
	rateLimiters = make(map[string]*rate.Limiter)
	mu           sync.RWMutex
)

func JWTAuth(secret string, adminSecret ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Debug: Log all headers for POST requests to /deposit
		if c.Request.Method == "POST" && strings.Contains(c.Request.URL.Path, "deposit") {
			log.Printf("🔐 [AUTH DEBUG] === Deposit Request ===")
			log.Printf("🔐 [AUTH DEBUG] Method: %s", c.Request.Method)
			log.Printf("🔐 [AUTH DEBUG] Path: %s", c.Request.URL.Path)
			log.Printf("🔐 [AUTH DEBUG] Remote IP: %s", c.ClientIP())
			log.Printf("🔐 [AUTH DEBUG] All Headers:")
			for key, values := range c.Request.Header {
				for _, value := range values {
					if key == "Authorization" {
						// Truncate token for security
						if len(value) > 30 {
							log.Printf("🔐 [AUTH DEBUG]   %s: %s...%s", key, value[:20], value[len(value)-10:])
						} else {
							log.Printf("🔐 [AUTH DEBUG]   %s: %s", key, value)
						}
					} else {
						log.Printf("🔐 [AUTH DEBUG]   %s: %s", key, value)
					}
				}
			}
		}

		tokenString := extractToken(c)

		// Debug: Log token extraction result
		if c.Request.Method == "POST" && strings.Contains(c.Request.URL.Path, "deposit") {
			if tokenString == "" {
				log.Printf("🔐 [AUTH DEBUG] ❌ Token extraction FAILED - no token found")
				log.Printf("🔐 [AUTH DEBUG] Raw Authorization header: '%s'", c.GetHeader("Authorization"))
			} else {
				log.Printf("🔐 [AUTH DEBUG] ✅ Token extracted successfully (length: %d)", len(tokenString))
			}
		}

		if tokenString == "" {
			log.Printf("🔐 [AUTH] No token for %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Try standard user token first
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		// Debug: Log parsing result for deposit requests
		if c.Request.Method == "POST" && strings.Contains(c.Request.URL.Path, "deposit") {
			if err != nil {
				log.Printf("🔐 [AUTH DEBUG] ❌ User token parse error: %v", err)
			} else {
				log.Printf("🔐 [AUTH DEBUG] ✅ User token parsed, valid: %v", token.Valid)
			}
		}

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(*Claims); ok {
				log.Printf("🔐 [AUTH] ✅ Authenticated user %s for %s %s", claims.UserID, c.Request.Method, c.Request.URL.Path)
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("user_role", claims.Role)
				c.Set("kyc_level", claims.KYCLevel)
				c.Next()
				return
			}
		}

		// If user token fails and we have an admin secret, try that
		if len(adminSecret) > 0 {
			// Parse as MapClaims since Admin token structure is different
			adminToken, adminErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(adminSecret[0]), nil
			})

			// Debug: Log admin token parsing
			if c.Request.Method == "POST" && strings.Contains(c.Request.URL.Path, "deposit") {
				if adminErr != nil {
					log.Printf("🔐 [AUTH DEBUG] ❌ Admin token parse error: %v", adminErr)
				} else {
					log.Printf("🔐 [AUTH DEBUG] Admin token parsed, valid: %v", adminToken.Valid)
				}
			}

			if adminErr == nil && adminToken.Valid {
				if claims, ok := adminToken.Claims.(jwt.MapClaims); ok {
					// Check if it's actually an admin token
					if claims["type"] == "admin" {
						log.Printf("🔐 [AUTH] ✅ Authenticated admin %s for %s %s", claims["admin_id"], c.Request.Method, c.Request.URL.Path)
						c.Set("user_id", claims["admin_id"])
						c.Set("user_email", claims["email"])
						c.Set("user_role", "admin") // Force admin role
						c.Next()
						return
					}
				}
			}
		}

		log.Printf("🔐 [AUTH] ❌ Invalid token for %s %s from %s (parse error: %v)", c.Request.Method, c.Request.URL.Path, c.ClientIP(), err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := rateLimiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(time.Minute/100), 100)
			rateLimiters[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")

	// Debug logging for token extraction
	if bearerToken == "" {
		log.Printf("🔐 [EXTRACT] Authorization header is empty")
		return ""
	}

	parts := strings.Split(bearerToken, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}

	log.Printf("🔐 [EXTRACT] Invalid Authorization format: '%s' (parts: %d)", bearerToken[:min(20, len(bearerToken))], len(parts))
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
