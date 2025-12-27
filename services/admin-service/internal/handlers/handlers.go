package handlers

import (
	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *services.AdminService
}

func NewAdminHandler(service *services.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

// ========== Authentication ==========

func (h *AdminHandler) Login(c *gin.Context) {
	var req models.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.service.CreateAuditLog(response.Admin.ID, response.Admin.Email, "login", "auth", "", "", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, response)
}

func (h *AdminHandler) GetCurrentAdmin(c *gin.Context) {
	adminID := c.GetString("admin_id")
	
	admin, err := h.service.GetAdmin(adminID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	permissions, _ := h.service.GetAdminPermissions(adminID)

	c.JSON(http.StatusOK, gin.H{
		"admin":       admin,
		"permissions": permissions,
	})
}

// ========== Admin CRUD ==========

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req models.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy := c.GetString("admin_id")
	admin, err := h.service.CreateAdmin(&req, createdBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "create_admin", "admins", admin.ID)
	c.JSON(http.StatusCreated, gin.H{"admin": admin})
}

func (h *AdminHandler) GetAdmin(c *gin.Context) {
	id := c.Param("id")
	
	admin, err := h.service.GetAdmin(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"admin": admin})
}

func (h *AdminHandler) GetAllAdmins(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	admins, err := h.service.GetAllAdmins(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"admins": admins})
}

func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	id := c.Param("id")
	
	var req models.UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateAdmin(id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "update_admin", "admins", id)
	c.JSON(http.StatusOK, gin.H{"message": "Admin updated"})
}

func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	id := c.Param("id")
	
	// Prevent self-deletion
	if id == c.GetString("admin_id") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete yourself"})
		return
	}

	if err := h.service.DeleteAdmin(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "delete_admin", "admins", id)
	c.JSON(http.StatusOK, gin.H{"message": "Admin deleted"})
}

// ========== Roles ==========

func (h *AdminHandler) GetRoles(c *gin.Context) {
	roles, err := h.service.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (h *AdminHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	
	role, err := h.service.GetRole(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

// ========== Dashboard ==========

func (h *AdminHandler) GetDashboard(c *gin.Context) {
	stats, err := h.service.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// ========== Users ==========

func (h *AdminHandler) GetUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := h.service.GetUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *AdminHandler) BlockUser(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetString("admin_id")
	if err := h.service.BlockUser(id, req.Reason, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "block_user", "users", id)
	c.JSON(http.StatusOK, gin.H{"message": "User blocked"})
}

func (h *AdminHandler) UnblockUser(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetString("admin_id")

	if err := h.service.UnblockUser(id, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "unblock_user", "users", id)
	c.JSON(http.StatusOK, gin.H{"message": "User unblocked"})
}

// ========== KYC ==========

func (h *AdminHandler) ApproveKYC(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Level string `json:"level" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetString("admin_id")
	if err := h.service.ApproveKYC(id, req.Level, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "approve_kyc", "kyc", id)
	c.JSON(http.StatusOK, gin.H{"message": "KYC approved"})
}

func (h *AdminHandler) RejectKYC(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetString("admin_id")
	if err := h.service.RejectKYC(id, req.Reason, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "reject_kyc", "kyc", id)
	c.JSON(http.StatusOK, gin.H{"message": "KYC rejected"})
}

func (h *AdminHandler) GetUserKYCDocuments(c *gin.Context) {
	userID := c.Param("id")

	docs, err := h.service.GetUserKYCDocuments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get KYC documents"})
		return
	}

	if docs == nil {
		docs = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, gin.H{"documents": docs})
}

// ========== Transactions ==========

func (h *AdminHandler) GetTransactions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transactions, err := h.service.GetTransactions(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func (h *AdminHandler) BlockTransaction(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetString("admin_id")
	if err := h.service.BlockTransaction(id, req.Reason, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "block_transaction", "transactions", id)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction blocked"})
}

func (h *AdminHandler) RefundTransaction(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetString("admin_id")
	if err := h.service.RefundTransaction(id, req.Reason, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "refund_transaction", "transactions", id)
	c.JSON(http.StatusOK, gin.H{"message": "Transaction refunded"})
}

// ========== Cards ==========

func (h *AdminHandler) GetCards(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	cards, err := h.service.GetCards(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cards": cards})
}

func (h *AdminHandler) FreezeCard(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetString("admin_id")
	if err := h.service.FreezeCard(id, req.Reason, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "freeze_card", "cards", id)
	c.JSON(http.StatusOK, gin.H{"message": "Card frozen"})
}

func (h *AdminHandler) BlockCard(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetString("admin_id")
	if err := h.service.BlockCard(id, req.Reason, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "block_card", "cards", id)
	c.JSON(http.StatusOK, gin.H{"message": "Card blocked"})
}

// ========== Wallets ==========

func (h *AdminHandler) GetWallets(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	wallets, err := h.service.GetWallets(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

func (h *AdminHandler) FreezeWallet(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetString("admin_id")
	if err := h.service.FreezeWallet(id, req.Reason, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logAction(c, "freeze_wallet", "wallets", id)
	c.JSON(http.StatusOK, gin.H{"message": "Wallet frozen"})
}

// ========== Audit Logs ==========

func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, err := h.service.GetAuditLogs(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// Helper
func (h *AdminHandler) logAction(c *gin.Context, action, resource, resourceID string) {
	adminID := c.GetString("admin_id")
	adminEmail := c.GetString("admin_email")
	h.service.CreateAuditLog(adminID, adminEmail, action, resource, resourceID, "", c.ClientIP(), c.Request.UserAgent())
}
