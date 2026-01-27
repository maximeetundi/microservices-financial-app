package handlers

import (
	"net/http"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/services"
	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	walletService *services.WalletService
	cryptoService *services.CryptoService
}

func NewAddressHandler(ws *services.WalletService, cs *services.CryptoService) *AddressHandler {
	return &AddressHandler{
		walletService: ws,
		cryptoService: cs,
	}
}

// GetDepositAddress godoc
// @Summary Get deposit address for a wallet
// @Description Get deposit address, optionally specifying a network (e.g., TRC20, ERC20)
// @Tags addresses
// @Produce json
// @Param id path string true "Wallet ID"
// @Param network query string false "Network (TRC20, ERC20, BEP20)"
// @Success 200 {object} map[string]string
// @Router /wallets/{id}/address [get]
func (h *AddressHandler) GetDepositAddress(c *gin.Context) {
	walletID := c.Param("id")
	network := c.DefaultQuery("network", "") // Optional network

	// 1. Get Wallet
	wallet, err := h.walletService.GetWalletByID(walletID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	// 2. Get User ID (In a real app, this would be from context/JWT)
	// For now, we trust the wallet's UserID
	userID := wallet.UserID

	// 3. Get or Create Address
	address, err := h.cryptoService.GetOrCreateAddress(userID, wallet.Currency, network)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate address: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address":  address,
		"network":  network,
		"currency": wallet.Currency,
	})
}
