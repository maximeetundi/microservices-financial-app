package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret string

func SetJWTSecret(secret string) {
	jwtSecret = secret
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Extract user info from token
		userID, _ := claims["sub"].(string)
		if userID == "" {
			userID, _ = claims["user_id"].(string)
		}
		
		email, _ := claims["email"].(string)
		firstName, _ := claims["first_name"].(string)
		lastName, _ := claims["last_name"].(string)

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("first_name", firstName)
		c.Set("last_name", lastName)
		c.Set("user_name", firstName + " " + lastName)

		c.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			return
		}

		userID, _ := claims["sub"].(string)
		if userID == "" {
			userID, _ = claims["user_id"].(string)
		}
		
		email, _ := claims["email"].(string)
		firstName, _ := claims["first_name"].(string)
		lastName, _ := claims["last_name"].(string)

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Set("first_name", firstName)
		c.Set("last_name", lastName)
		c.Set("user_name", firstName + " " + lastName)

		c.Next()
	}
}
