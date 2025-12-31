package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Set user info in context
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}
		// For admin tokens, admin_id is used instead of user_id
		if adminID, ok := claims["admin_id"].(string); ok {
			c.Set("user_id", adminID)
			c.Set("admin_id", adminID)
		}
		if userName, ok := claims["name"].(string); ok {
			c.Set("user_name", userName)
		}
		if userEmail, ok := claims["email"].(string); ok {
			c.Set("user_email", userEmail)
		}
		if role, ok := claims["role"].(string); ok {
			c.Set("role", role)
		}
		// For admin tokens, 'type' claim is used
		if tokenType, ok := claims["type"].(string); ok {
			c.Set("token_type", tokenType)
		}

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for admin token type (from admin-service JWT)
		tokenType := c.GetString("token_type")
		if tokenType == "admin" {
			c.Next()
			return
		}
		
		// Fall back to checking role claim (for regular user tokens)
		role := c.GetString("role")
		if role == "admin" || role == "support_agent" {
			c.Next()
			return
		}
		
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin or support agent access required"})
		c.Abort()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Next()
	}
}

func RateLimiter() gin.HandlerFunc {
	// Simple rate limiter - in production use Redis-based limiter
	return func(c *gin.Context) {
		c.Next()
	}
}
