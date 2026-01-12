package handlers

	"net/http"
	"strconv"

	"github.com/crypto-bank/microservices-financial-app/services/donation-service/internal/services"
	"github.com/gin-gonic/gin"
)

type DonationHandler struct {
	service *services.DonationService
}

func NewDonationHandler(service *services.DonationService) *DonationHandler {
	return &DonationHandler{service: service}
}

func (h *DonationHandler) Initiate(c *gin.Context) {
	var req services.InitiateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetHeader("X-User-ID")
	if userID != "" {
		req.DonorID = userID
	} else {
		// If no user ID, is it anonymous guest? For now assuming logic requires ID or empty is handled.
		// "Anonymous" flag means hide name, but maybe we still track who paid if logged in?
		// If guest, DonorID is empty.
	}

	donation, err := h.service.InitiateDonation(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "donation initiated", "donation": donation})
}

func (h *DonationHandler) List(c *gin.Context) {
	campaignID := c.Query("campaign_id")
	if campaignID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "campaign_id required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	donations, err := h.service.ListDonations(campaignID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Post-process to hide info if anonymous? 
	// Or define a public response model.
	// For MVP, Frontend will hide it. Or we zero it out here.
	for _, d := range donations {
		if d.IsAnonymous {
			d.DonorID = "" // Hide sensitive ID if anonymous
			// Also FormData might contain name, hide it?
			d.FormData = nil
		}
	}

	c.JSON(http.StatusOK, gin.H{"donations": donations})
}
