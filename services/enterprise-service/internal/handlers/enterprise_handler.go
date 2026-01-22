package handlers

import (
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/repository"
	"github.com/crypto-bank/microservices-financial-app/services/enterprise-service/internal/services"
	"github.com/gin-gonic/gin"
)

type EnterpriseHandler struct {
	repo    *repository.EnterpriseRepository
	empRepo *repository.EmployeeRepository
	storage *services.StorageService
}

func NewEnterpriseHandler(repo *repository.EnterpriseRepository, empRepo *repository.EmployeeRepository, storage *services.StorageService) *EnterpriseHandler {
	return &EnterpriseHandler{repo: repo, empRepo: empRepo, storage: storage}
}

// UploadLogo handles logo upload
func (h *EnterpriseHandler) UploadLogo(c *gin.Context) {
	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	// Upload to MinIO (folder "enterprises")
	url, err := h.storage.UploadFile(c.Request.Context(), file, header.Filename, header.Size, contentType, "enterprises")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// CreateEnterprise (Point 1)
func (h *EnterpriseHandler) CreateEnterprise(c *gin.Context) {
	var req models.Enterprise
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set OwnerID from JWT (assuming middleware sets "user_id")
	userID, _ := c.Get("user_id")
	req.OwnerID, _ = userID.(string)

	// Set Default Type if missing
	if req.Type == "" {
		req.Type = models.EnterpriseTypeService
	}

	// Add default security policies if none provided
	if len(req.SecurityPolicies) == 0 {
		req.SecurityPolicies = getDefaultSecurityPolicies()
	}

	if err := h.repo.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create enterprise", "details": err.Error()})
		return
	}

	// Auto-create owner as admin employee with full permissions
	if h.empRepo != nil && req.OwnerID != "" {
		ownerEmployee := &models.Employee{
			EnterpriseID: req.ID,
			UserID:       req.OwnerID,
			FirstName:    "Administrateur",
			LastName:     "Principal",
			Profession:   "Propriétaire",
			Role:         models.EmployeeRoleOwner, // Owner role = all permissions
			Status:       models.EmployeeStatusActive,
			Permissions:  models.GetDefaultOwnerPermissions(), // Full permissions
			HireDate:     time.Now(),
			InvitedAt:    time.Now(),
			AcceptedAt:   time.Now(),
		}
		// Ignore error - enterprise creation should still succeed
		_ = h.empRepo.Create(c.Request.Context(), ownerEmployee)
	}

	c.JSON(http.StatusCreated, req)
}

// getDefaultSecurityPolicies returns default security policies for a new enterprise
func getDefaultSecurityPolicies() []models.SecurityPolicy {
	return []models.SecurityPolicy{
		{
			ID:               "default_transaction",
			Name:             "Transactions importantes",
			ActionType:       models.ActionTypeTransaction,
			Enabled:          true,
			MinApprovals:     1,
			RequireMajority:  false,
			RequireAllAdmins: true,
			ThresholdAmount:  100000, // 100,000 XOF threshold
			ExpirationHours:  24,
		},
		{
			ID:               "default_payroll",
			Name:             "Paiement des salaires",
			ActionType:       models.ActionTypePayroll,
			Enabled:          true,
			MinApprovals:     1,
			RequireMajority:  false,
			RequireAllAdmins: true,
			ThresholdAmount:  0, // Always requires approval
			ExpirationHours:  48,
		},
		{
			ID:               "default_admin_change",
			Name:             "Modification des administrateurs",
			ActionType:       models.ActionTypeAdminChange,
			Enabled:          true,
			MinApprovals:     1,
			RequireMajority:  false,
			RequireAllAdmins: true,
			ThresholdAmount:  0,
			ExpirationHours:  24,
		},
		{
			ID:               "default_settings",
			Name:             "Modification des paramètres",
			ActionType:       models.ActionTypeSettingsChange,
			Enabled:          true,
			MinApprovals:     1,
			RequireMajority:  false,
			RequireAllAdmins: true,
			ThresholdAmount:  0,
			ExpirationHours:  24,
		},
		{
			ID:               "default_employee_terminate",
			Name:             "Licenciement d'employé",
			ActionType:       models.ActionTypeEmployeeTerminate,
			Enabled:          true,
			MinApprovals:     1,
			RequireMajority:  true, // Majority for hiring/firing
			RequireAllAdmins: false,
			ThresholdAmount:  0,
			ExpirationHours:  48,
		},
		{
			ID:               "default_enterprise_delete",
			Name:             "Suppression de l'entreprise",
			ActionType:       models.ActionTypeEnterpriseDelete,
			Enabled:          true,
			MinApprovals:     1,
			RequireMajority:  false,
			RequireAllAdmins: true, // ALL admins must approve deletion
			ThresholdAmount:  0,
			ExpirationHours:  72, // 3 days to decide
		},
	}
}

func (h *EnterpriseHandler) GetEnterprise(c *gin.Context) {
	id := c.Param("id")
	ent, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enterprise not found"})
		return
	}
	c.JSON(http.StatusOK, ent)
}

func (h *EnterpriseHandler) ListEnterprises(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Get enterprises where user is owner OR is an active employee
	enterprises, err := h.repo.FindByUserAccess(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list enterprises", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enterprises)
}

// UpdateEnterprise (Required for Settings)
func (h *EnterpriseHandler) UpdateEnterprise(c *gin.Context) {
	id := c.Param("id")
	
	// Fetch existing to ensure ownership & existence
	existing, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enterprise not found"})
		return
	}

	// Verify Owner (simple check)
	userID, exists := c.Get("user_id")
	if !exists || existing.OwnerID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update this enterprise"})
		return
	}

	// Bind updates
	var req models.Enterprise
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve ID and Owner
	req.ID = existing.ID
	req.OwnerID = existing.OwnerID
	req.CreatedAt = existing.CreatedAt // Don't overwrite creation time

	if err := h.repo.Update(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update enterprise", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}



// GetEnterpriseNotifications fetches notifications for the enterprise
func (h *EnterpriseHandler) GetEnterpriseNotifications(c *gin.Context) {
	id := c.Param("id") // Enterprise ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify User Access (Must be an employee of the enterprise)
	// We can check if they are an employee.
	isEmployee, err := h.empRepo.IsEmployee(c.Request.Context(), id, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify permissions"})
		return
	}
	if !isEmployee {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Proxy to Notification Service
	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationServiceURL == "" {
		notificationServiceURL = "http://notification-service:8087"
	}

	resp, err := http.Get(notificationServiceURL + "/api/v1/internal/users/" + id + "/notifications?limit=50")
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Notification service unavailable"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	// Stream response back
	c.DataFromReader(http.StatusOK, resp.ContentLength, "application/json", resp.Body, nil)
}

// ========== ADMIN ENDPOINTS ==========

// AdminListAll returns all enterprises for admin dashboard
// GET /admin/enterprises
func (h *EnterpriseHandler) AdminListAll(c *gin.Context) {
	enterprises, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enterprises"})
		return
	}
	
	// Add employee counts
	result := make([]map[string]interface{}, len(enterprises))
	for i, ent := range enterprises {
		empCount, _ := h.empRepo.CountByEnterprise(c.Request.Context(), ent.ID.Hex())
		result[i] = map[string]interface{}{
			"id":              ent.ID.Hex(),
			"name":            ent.Name,
			"type":            ent.Type,
			"status":          ent.Status,
			"owner_id":        ent.OwnerID,
			"logo":            ent.Logo,
			"description":     ent.Description,
			"employees_count": empCount,
			"created_at":      ent.CreatedAt,
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"enterprises": result})
}

// AdminUpdateStatus changes an enterprise status (activate/suspend)
// PUT /admin/enterprises/:id/status
func (h *EnterpriseHandler) AdminUpdateStatus(c *gin.Context) {
	id := c.Param("id")
	
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate status
	validStatuses := []string{"ACTIVE", "SUSPENDED", "PENDING"}
	valid := false
	for _, s := range validStatuses {
		if req.Status == s {
			valid = true
			break
		}
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be ACTIVE, SUSPENDED, or PENDING"})
		return
	}
	
	// Get enterprise
	enterprise, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enterprise not found"})
		return
	}
	
	// Update status
	enterprise.Status = req.Status
	if err := h.repo.Update(c.Request.Context(), enterprise); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Status updated", "status": req.Status})
}

// AdminDelete deletes an enterprise and all related data
// DELETE /admin/enterprises/:id
func (h *EnterpriseHandler) AdminDelete(c *gin.Context) {
	id := c.Param("id")
	
	// Get enterprise first to confirm it exists
	_, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enterprise not found"})
		return
	}
	
	// Delete employees first
	if h.empRepo != nil {
		h.empRepo.DeleteByEnterprise(c.Request.Context(), id)
	}
	
	// Delete enterprise
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete enterprise"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Enterprise deleted"})
}

