package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AdminAuth middleware for admin authentication
func AdminAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Check if it's an admin token
		if claims["type"] != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not an admin token"})
			return
		}

		// Set admin info in context
		c.Set("admin_id", claims["admin_id"])
		c.Set("admin_email", claims["email"])
		c.Set("admin_role_id", claims["role_id"])
		
		// Set permissions
		if permissions, ok := claims["permissions"].([]interface{}); ok {
			permStrings := make([]string, len(permissions))
			for i, p := range permissions {
				permStrings[i] = p.(string)
			}
			c.Set("admin_permissions", permStrings)
		}

		c.Next()
	}
}

// RequirePermission middleware to check specific permission
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("admin_permissions")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
			return
		}

		perms, ok := permissions.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
			return
		}

		hasPermission := false
		for _, p := range perms {
			if p == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":      "Permission denied",
				"required":   permission,
			})
			return
		}

		c.Next()
	}
}

// RequireAnyPermission middleware to check if user has any of the specified permissions
func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminPerms, exists := c.Get("admin_permissions")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No permissions found"})
			return
		}

		perms, ok := adminPerms.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid permissions format"})
			return
		}

		hasPermission := false
		for _, required := range permissions {
			for _, has := range perms {
				if required == has {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":    "Permission denied",
				"required": permissions,
			})
			return
		}

		c.Next()
	}
}

// SecurityHeaders adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}

// AuditLog middleware for logging admin actions
type AuditLogger func(adminID, adminEmail, action, resource, resourceID, details, ip, userAgent string)

func AuditMiddleware(logger AuditLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only log successful mutating requests
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			method := c.Request.Method
			if method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE" {
				adminID, _ := c.Get("admin_id")
				adminEmail, _ := c.Get("admin_email")
				
				logger(
					adminID.(string),
					adminEmail.(string),
					method,
					c.FullPath(),
					c.Param("id"),
					"",
					c.ClientIP(),
					c.Request.UserAgent(),
				)
			}
		}
	}
}

// CORS middleware for admin dashboard
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
