package middleware

import (
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
		tokenString := extractToken(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Try standard user token first
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(*Claims); ok {
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

			if adminErr == nil && adminToken.Valid {
				if claims, ok := adminToken.Claims.(jwt.MapClaims); ok {
					// Check if it's actually an admin token
					if claims["type"] == "admin" {
						c.Set("user_id", claims["admin_id"])
						c.Set("user_email", claims["email"])
						c.Set("user_role", "admin") // Force admin role
						c.Next()
						return
					}
				}
			}
		}

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
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
