package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/middleware"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/repository"
	"github.com/gin-gonic/gin"
)

// ExtendedHandler handles the new features (roles, approvals, chat, solidarity)
type ExtendedHandler struct {
	DB            *sql.DB // For emergency fund operations
	roleRepo      *repository.RoleRepository
	approvalRepo  *repository.ApprovalRepository
	chatRepo      *repository.ChatRepository
	solidarityRepo *repository.SolidarityRepository
	calledRepo    *repository.CalledRoundRepository
	memberRepo    *repository.MemberRepository
}

func NewExtendedHandler(
	db *sql.DB,
	roleRepo *repository.RoleRepository,
	approvalRepo *repository.ApprovalRepository,
	chatRepo *repository.ChatRepository,
	solidarityRepo *repository.SolidarityRepository,
	calledRepo *repository.CalledRoundRepository,
	memberRepo *repository.MemberRepository,
) *ExtendedHandler {
	return &ExtendedHandler{
		DB:            db,
		roleRepo:      roleRepo,
		approvalRepo:  approvalRepo,
		chatRepo:      chatRepo,
		solidarityRepo: solidarityRepo,
		calledRepo:    calledRepo,
		memberRepo:    memberRepo,
	}
}

// === Custom Roles ===

func (h *ExtendedHandler) CreateRole(c *gin.Context) {
	associationID := c.Param("id")

	var req models.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := &models.AssociationRole{
		AssociationID: associationID,
		Name:          req.Name,
		Permissions:   req.Permissions,
		IsDefault:     false,
	}

	if err := h.roleRepo.Create(role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

func (h *ExtendedHandler) GetRoles(c *gin.Context) {
	associationID := c.Param("id")

	roles, err := h.roleRepo.GetByAssociation(associationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (h *ExtendedHandler) DeleteRole(c *gin.Context) {
	roleID := c.Param("roleId")

	if err := h.roleRepo.Delete(roleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}

// === Multi-Signature Approvers ===

func (h *ExtendedHandler) SetApprovers(c *gin.Context) {
	associationID := c.Param("id")

	var req models.SetApproversRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.MemberIDs) != 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exactly 5 approvers required"})
		return
	}

	if err := h.approvalRepo.SetApprovers(associationID, req.MemberIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Approvers set"})
}

func (h *ExtendedHandler) GetApprovers(c *gin.Context) {
	associationID := c.Param("id")

	approvers, err := h.approvalRepo.GetApprovers(associationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, approvers)
}

func (h *ExtendedHandler) GetPendingApprovals(c *gin.Context) {
	associationID := c.Param("id")

	requests, err := h.approvalRepo.GetPendingRequests(associationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Populate vote counts
	for _, req := range requests {
		approvals, rejections, _ := h.approvalRepo.CountVotes(req.ID)
		req.CurrentApprovals = approvals
		req.CurrentRejections = rejections
	}

	c.JSON(http.StatusOK, requests)
}

func (h *ExtendedHandler) VoteOnApproval(c *gin.Context) {
	requestID := c.Param("requestId")
	member, exists := c.Get("member")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not a member"})
		return
	}
	memberObj := member.(*models.Member)

	var req models.VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if this member is an approver
	approvalReq, err := h.approvalRepo.GetRequest(requestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	isApprover, _ := h.approvalRepo.IsApprover(approvalReq.AssociationID, memberObj.ID)
	if !isApprover {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not an approver"})
		return
	}

	vote := &models.ApprovalVote{
		RequestID:  requestID,
		ApproverID: memberObj.ID,
		Vote:       req.Vote,
		Comment:    req.Comment,
	}

	if err := h.approvalRepo.Vote(vote); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if we have enough votes
	approvals, rejections, _ := h.approvalRepo.CountVotes(requestID)
	
	if approvals >= approvalReq.RequiredApprovals {
		h.approvalRepo.UpdateRequestStatus(requestID, models.ApprovalStatusApproved)
		// TODO: Execute the approved action (loan disbursement, distribution, etc.)
	} else if rejections >= 2 {
		h.approvalRepo.UpdateRequestStatus(requestID, models.ApprovalStatusRejected)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vote recorded",
		"approvals": approvals,
		"rejections": rejections,
		"required": approvalReq.RequiredApprovals,
	})
}

// === Chat ===

func (h *ExtendedHandler) SendMessage(c *gin.Context) {
	associationID := c.Param("id")
	member, exists := c.Get("member")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not a member"})
		return
	}
	memberObj := member.(*models.Member)

	var req models.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := &models.ChatMessage{
		AssociationID: associationID,
		SenderID:      memberObj.ID,
		Content:       req.Content,
		IsAdminOnly:   req.IsAdminOnly,
	}

	if err := h.chatRepo.Create(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

func (h *ExtendedHandler) GetMessages(c *gin.Context) {
	associationID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	// Check if user is admin to see admin-only messages
	isAdmin := false
	member, exists := c.Get("member")
	if exists {
		memberObj := member.(*models.Member)
		isAdmin = memberObj.Role == models.MemberRolePresident || 
				  memberObj.Role == models.MemberRoleTreasurer ||
				  memberObj.Role == models.MemberRoleSecretary
	}

	messages, err := h.chatRepo.GetByAssociation(associationID, limit, offset, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// === Solidarity Events ===

func (h *ExtendedHandler) CreateSolidarityEvent(c *gin.Context) {
	userID := middleware.GetUserID(c)
	associationID := c.Param("id")

	var req models.CreateSolidarityEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &models.SolidarityEvent{
		AssociationID: associationID,
		EventType:     req.EventType,
		BeneficiaryID: req.BeneficiaryID,
		Title:         req.Title,
		Description:   req.Description,
		TargetAmount:  req.TargetAmount,
		Status:        models.SolidarityStatusActive,
		CreatedBy:     userID,
	}

	if err := h.solidarityRepo.CreateEvent(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *ExtendedHandler) GetSolidarityEvents(c *gin.Context) {
	associationID := c.Param("id")
	status := c.Query("status")

	events, err := h.solidarityRepo.GetByAssociation(associationID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *ExtendedHandler) ContributeToSolidarity(c *gin.Context) {
	eventID := c.Param("eventId")
	member, exists := c.Get("member")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not a member"})
		return
	}
	memberObj := member.(*models.Member)

	var req models.ContributeSolidarityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contrib := &models.SolidarityContribution{
		EventID:       eventID,
		ContributorID: memberObj.ID,
		Amount:        req.Amount,
		Paid:          true, // Assuming payment via wallet is immediate
	}

	if err := h.solidarityRepo.AddContribution(contrib); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update collected amount
	h.solidarityRepo.UpdateCollectedAmount(eventID, req.Amount)

	c.JSON(http.StatusCreated, contrib)
}

// === Called Tontine ===

func (h *ExtendedHandler) CreateCalledRound(c *gin.Context) {
	associationID := c.Param("id")

	var req models.CreateCalledRoundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roundNumber, _ := h.calledRepo.GetNextRoundNumber(associationID)

	round := &models.CalledRound{
		AssociationID: associationID,
		BeneficiaryID: req.BeneficiaryID,
		RoundNumber:   roundNumber,
		Status:        "active",
	}

	if err := h.calledRepo.Create(round); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, round)
}

func (h *ExtendedHandler) GetCalledRounds(c *gin.Context) {
	associationID := c.Param("id")
	status := c.Query("status")

	rounds, err := h.calledRepo.GetByAssociation(associationID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rounds)
}

func (h *ExtendedHandler) MakePledge(c *gin.Context) {
	roundID := c.Param("roundId")
	member, exists := c.Get("member")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not a member"})
		return
	}
	memberObj := member.(*models.Member)

	var req models.MakePledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pledge := &models.CalledPledge{
		RoundID:       roundID,
		ContributorID: memberObj.ID,
		Amount:        req.Amount,
	}

	if err := h.calledRepo.AddPledge(pledge); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pledge)
}

func (h *ExtendedHandler) PayPledge(c *gin.Context) {
	pledgeID := c.Param("pledgeId")
	roundID := c.Param("roundId")

	var req models.PayPledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Integrate wallet service for actual payment

	if err := h.calledRepo.MarkPledgePaid(pledgeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get pledge amount to update total
	pledges, _ := h.calledRepo.GetPledges(roundID)
	for _, p := range pledges {
		if p.ID == pledgeID {
			h.calledRepo.UpdateTotalCollected(roundID, p.Amount)
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pledge paid"})
}

func (h *ExtendedHandler) GetPledges(c *gin.Context) {
	roundID := c.Param("roundId")

	pledges, err := h.calledRepo.GetPledges(roundID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pledges)
}
