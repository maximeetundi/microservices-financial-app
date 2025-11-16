package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crypto-bank/exchange-service/internal/models"
	"github.com/crypto-bank/exchange-service/internal/services"
)

type FiatHandler struct {
	fiatExchangeService *services.FiatExchangeService
	rateService         *services.RateService
}

func NewFiatHandler(fiatExchangeService *services.FiatExchangeService, rateService *services.RateService) *FiatHandler {
	return &FiatHandler{
		fiatExchangeService: fiatExchangeService,
		rateService:         rateService,
	}
}

// Obtenir un devis pour conversion fiat
func (h *FiatHandler) GetFiatQuote(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		FromCurrency string  `json:"from_currency" binding:"required"`
		ToCurrency   string  `json:"to_currency" binding:"required"`
		Amount       float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quote, err := h.fiatExchangeService.GetFiatQuote(
		userID.(string),
		req.FromCurrency,
		req.ToCurrency,
		req.Amount,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quote": quote,
		"exchange_type": "fiat_to_fiat",
		"estimated_delivery": "Instant",
	})
}

// Exécuter une conversion fiat
func (h *FiatHandler) ExecuteFiatExchange(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		FromWalletID string  `json:"from_wallet_id" binding:"required"`
		ToWalletID   string  `json:"to_wallet_id" binding:"required"`
		FromCurrency string  `json:"from_currency" binding:"required"`
		ToCurrency   string  `json:"to_currency" binding:"required"`
		Amount       float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchange, err := h.fiatExchangeService.ConvertFiat(
		userID.(string),
		req.FromCurrency,
		req.ToCurrency,
		req.Amount,
		req.FromWalletID,
		req.ToWalletID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Fiat exchange completed successfully",
		"exchange": exchange,
		"processing_time": "Instant",
	})
}

// Obtenir les taux fiat en temps réel
func (h *FiatHandler) GetFiatRates(c *gin.Context) {
	baseCurrency := c.DefaultQuery("base", "USD")
	
	// Obtenir tous les taux par rapport à la devise de base
	supportedPairs := h.fiatExchangeService.GetSupportedFiatPairs()
	targetCurrencies, exists := supportedPairs[baseCurrency]
	
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported base currency",
			"supported_currencies": []string{"USD", "EUR", "GBP", "JPY", "CAD", "AUD"},
		})
		return
	}

	rates := make(map[string]interface{})
	
	for _, targetCurrency := range targetCurrencies {
		rate, err := h.rateService.GetRate(baseCurrency, targetCurrency)
		if err == nil {
			rates[targetCurrency] = map[string]interface{}{
				"rate":        rate.Rate,
				"bid":         rate.BidPrice,
				"ask":         rate.AskPrice,
				"change_24h":  rate.Change24h,
				"last_update": rate.LastUpdated,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"base_currency": baseCurrency,
		"rates": rates,
		"timestamp": "real_time",
		"source": "fiat_exchange_api",
	})
}

// Obtenir un taux spécifique fiat
func (h *FiatHandler) GetSpecificFiatRate(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")

	rate, err := h.rateService.GetRate(from, to)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exchange rate not found"})
		return
	}

	// Calculer les frais pour cette conversion
	sampleAmount := 1000.0 // Montant de référence
	feePercentage := h.fiatExchangeService.CalculateFiatExchangeFee(from, to, sampleAmount)

	c.JSON(http.StatusOK, gin.H{
		"pair": from + "/" + to,
		"rate": rate.Rate,
		"bid": rate.BidPrice,
		"ask": rate.AskPrice,
		"spread": rate.Spread,
		"change_24h": rate.Change24h,
		"volume_24h": rate.Volume24h,
		"fee_percentage": feePercentage,
		"last_updated": rate.LastUpdated,
		"is_fiat_pair": true,
	})
}

// Obtenir l'historique des taux
func (h *FiatHandler) GetFiatRateHistory(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")
	period := c.DefaultQuery("period", "7d") // 1d, 7d, 30d, 1y
	
	// TODO: Implémenter l'historique des taux
	// Pour l'instant, retourner des données simulées
	
	history := []map[string]interface{}{}
	
	// Générer un historique simulé
	days := 7
	if period == "1d" {
		days = 1
	} else if period == "30d" {
		days = 30
	} else if period == "1y" {
		days = 365
	}

	// Obtenir le taux actuel comme base
	currentRate, err := h.rateService.GetRate(from, to)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rate not found"})
		return
	}

	baseRate := currentRate.Rate
	
	for i := days; i >= 0; i-- {
		// Simuler une variation de +/- 2%
		variation := 1.0 + (0.04*(float64(i%10)/10.0) - 0.02)
		historicalRate := baseRate * variation
		
		history = append(history, map[string]interface{}{
			"date":  "2024-01-" + strconv.Itoa(30-i),
			"rate":  historicalRate,
			"high":  historicalRate * 1.01,
			"low":   historicalRate * 0.99,
			"close": historicalRate,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"pair": from + "/" + to,
		"period": period,
		"history": history,
		"current_rate": currentRate.Rate,
	})
}

// Calculatrice de conversion fiat
func (h *FiatHandler) FiatConverter(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")
	
	if from == "" || to == "" || amountStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing parameters",
			"required": []string{"from", "to", "amount"},
			"example": "/fiat/convert?from=USD&to=EUR&amount=100",
		})
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	// Obtenir le taux de change
	rate, err := h.rateService.GetRate(from, to)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exchange rate not found"})
		return
	}

	// Calculer la conversion
	convertedAmount := amount * rate.Rate
	
	// Calculer les frais
	feePercentage := h.fiatExchangeService.CalculateFiatExchangeFee(from, to, amount)
	fee := amount * feePercentage / 100
	finalAmount := convertedAmount - (fee * rate.Rate)

	c.JSON(http.StatusOK, gin.H{
		"from_currency": from,
		"to_currency": to,
		"original_amount": amount,
		"exchange_rate": rate.Rate,
		"converted_amount": convertedAmount,
		"fee": fee,
		"fee_percentage": feePercentage,
		"final_amount": finalAmount,
		"savings_vs_banks": map[string]interface{}{
			"traditional_bank_fee": amount * 0.035, // 3.5% frais bancaires typiques
			"our_fee": fee,
			"you_save": (amount * 0.035) - fee,
		},
		"breakdown": map[string]interface{}{
			"base_conversion": convertedAmount,
			"our_fee": fee,
			"net_amount": finalAmount,
		},
	})
}

// Obtenir les devises supportées pour conversions fiat
func (h *FiatHandler) GetSupportedFiatCurrencies(c *gin.Context) {
	currencies := map[string]map[string]interface{}{
		"USD": {
			"name": "US Dollar",
			"symbol": "$",
			"flag": "🇺🇸",
			"popular": true,
		},
		"EUR": {
			"name": "Euro",
			"symbol": "€", 
			"flag": "🇪🇺",
			"popular": true,
		},
		"GBP": {
			"name": "British Pound",
			"symbol": "£",
			"flag": "🇬🇧",
			"popular": true,
		},
		"JPY": {
			"name": "Japanese Yen",
			"symbol": "¥",
			"flag": "🇯🇵",
			"popular": true,
		},
		"CAD": {
			"name": "Canadian Dollar",
			"symbol": "C$",
			"flag": "🇨🇦",
			"popular": true,
		},
		"AUD": {
			"name": "Australian Dollar",
			"symbol": "A$",
			"flag": "🇦🇺",
			"popular": true,
		},
		"CHF": {
			"name": "Swiss Franc",
			"symbol": "CHF",
			"flag": "🇨🇭",
			"popular": false,
		},
		"SEK": {
			"name": "Swedish Krona",
			"symbol": "SEK",
			"flag": "🇸🇪",
			"popular": false,
		},
		"NOK": {
			"name": "Norwegian Krone",
			"symbol": "NOK",
			"flag": "🇳🇴",
			"popular": false,
		},
		"DKK": {
			"name": "Danish Krone",
			"symbol": "DKK",
			"flag": "🇩🇰",
			"popular": false,
		},
	}

	supportedPairs := h.fiatExchangeService.GetSupportedFiatPairs()

	c.JSON(http.StatusOK, gin.H{
		"currencies": currencies,
		"supported_pairs": supportedPairs,
		"total_currencies": len(currencies),
		"total_pairs": len(supportedPairs),
		"features": []string{
			"Real-time rates",
			"Low fees (0.15-0.25%)",
			"Instant transfers", 
			"24/7 availability",
			"No hidden charges",
		},
	})
}

// Comparer avec les frais bancaires traditionnels
func (h *FiatHandler) CompareBankingFees(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	amountStr := c.Query("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		amount = 1000.0 // Montant par défaut
	}

	if from == "" {
		from = "USD"
	}
	if to == "" {
		to = "EUR"
	}

	// Nos frais
	ourFeePercentage := h.fiatExchangeService.CalculateFiatExchangeFee(from, to, amount)
	ourFee := amount * ourFeePercentage / 100

	// Frais bancaires traditionnels (typiques)
	traditionalFees := map[string]interface{}{
		"big_banks": map[string]interface{}{
			"fee_percentage": 3.5,
			"fixed_fee": 15.0,
			"total_fee": (amount * 0.035) + 15.0,
			"exchange_markup": 4.0, // 4% markup sur le taux
		},
		"online_banks": map[string]interface{}{
			"fee_percentage": 2.0,
			"fixed_fee": 5.0,
			"total_fee": (amount * 0.02) + 5.0,
			"exchange_markup": 2.5,
		},
		"money_transfer": map[string]interface{}{
			"fee_percentage": 1.5,
			"fixed_fee": 10.0,
			"total_fee": (amount * 0.015) + 10.0,
			"exchange_markup": 2.0,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"amount": amount,
		"currency_pair": from + " → " + to,
		"crypto_bank": map[string]interface{}{
			"fee_percentage": ourFeePercentage,
			"fee_amount": ourFee,
			"exchange_markup": 0.1, // Notre markup minimal
			"delivery_time": "Instant",
		},
		"traditional_providers": traditionalFees,
		"savings": map[string]interface{}{
			"vs_big_banks": ((amount * 0.035) + 15.0) - ourFee,
			"vs_online_banks": ((amount * 0.02) + 5.0) - ourFee,
			"vs_money_transfer": ((amount * 0.015) + 10.0) - ourFee,
		},
		"advantages": []string{
			"No hidden fees",
			"Real-time exchange rates", 
			"Instant transfers",
			"24/7 availability",
			"Multi-currency support",
			"Crypto integration",
		},
	})
}