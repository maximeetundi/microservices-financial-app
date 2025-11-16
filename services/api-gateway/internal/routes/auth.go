package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/api-gateway/internal/services"
)

func SetupAuthRoutes(router *gin.RouterGroup, serviceManager *services.ServiceManager) {
	router.POST("/register", handleRegister(serviceManager))
	router.POST("/login", handleLogin(serviceManager))
	router.POST("/refresh", handleRefreshToken(serviceManager))
	router.POST("/logout", handleLogout(serviceManager))
	router.POST("/forgot-password", handleForgotPassword(serviceManager))
	router.POST("/reset-password", handleResetPassword(serviceManager))
	router.POST("/verify-email", handleVerifyEmail(serviceManager))
	router.POST("/verify-phone", handleVerifyPhone(serviceManager))
	router.POST("/enable-2fa", handleEnable2FA(serviceManager))
	router.POST("/verify-2fa", handleVerify2FA(serviceManager))
}

func handleRegister(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email       string `json:"email" binding:"required,email"`
			Phone       string `json:"phone" binding:"required"`
			Password    string `json:"password" binding:"required,min=8"`
			FirstName   string `json:"first_name" binding:"required"`
			LastName    string `json:"last_name" binding:"required"`
			Country     string `json:"country" binding:"required,len=3"`
			DateOfBirth string `json:"date_of_birth" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call auth service
		resp, err := sm.Register(c.Request.Context(), map[string]interface{}{
			"email":         req.Email,
			"phone":         req.Phone,
			"password":      req.Password,
			"first_name":    req.FirstName,
			"last_name":     req.LastName,
			"country":       req.Country,
			"date_of_birth": req.DateOfBirth,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleLogin(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
			TwoFACode string `json:"two_fa_code"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call auth service
		resp, err := sm.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleRefreshToken(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := sm.RefreshToken(c.Request.Context(), req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleLogout(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token and call auth service to invalidate
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
			return
		}

		// Call auth service to logout
		resp, err := sm.CallService(c.Request.Context(), "auth", "POST", "/logout", nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleForgotPassword(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := sm.CallService(c.Request.Context(), "auth", "POST", "/forgot-password", req, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleResetPassword(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Token       string `json:"token" binding:"required"`
			NewPassword string `json:"new_password" binding:"required,min=8"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := sm.CallService(c.Request.Context(), "auth", "POST", "/reset-password", req, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleVerifyEmail(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Token string `json:"token" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := sm.CallService(c.Request.Context(), "auth", "POST", "/verify-email", req, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleVerifyPhone(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Phone string `json:"phone" binding:"required"`
			Code  string `json:"code" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := sm.CallService(c.Request.Context(), "auth", "POST", "/verify-phone", req, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleEnable2FA(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		
		resp, err := sm.CallService(c.Request.Context(), "auth", "POST", "/enable-2fa", nil, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}

func handleVerify2FA(sm *services.ServiceManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Code string `json:"code" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token := c.GetHeader("Authorization")
		
		resp, err := sm.CallService(c.Request.Context(), "auth", "POST", "/verify-2fa", req, map[string]string{
			"Authorization": token,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service unavailable"})
			return
		}

		c.Data(resp.StatusCode, "application/json", resp.Body)
	}
}