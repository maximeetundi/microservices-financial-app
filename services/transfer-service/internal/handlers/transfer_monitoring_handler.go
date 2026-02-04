package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TransferMonitoringHandler handles monitoring endpoints for transfers
type TransferMonitoringHandler struct {
	db *sql.DB
}

func NewTransferMonitoringHandler(db *sql.DB) *TransferMonitoringHandler {
	return &TransferMonitoringHandler{db: db}
}

// GetMetrics returns real-time transfer metrics
// GET /api/v1/admin/transfers/metrics
func (h *TransferMonitoringHandler) GetMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	today := time.Now().Truncate(24 * time.Hour)

	metrics := struct {
		TodayDeposits    MetricData `json:"today_deposits"`
		TodayWithdrawals MetricData `json:"today_withdrawals"`
		Pending          MetricData `json:"pending"`
		Failed           MetricData `json:"failed"`
	}{}

	// Today's deposits
	err := h.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE type = 'deposit'
		AND created_at >= $1
		AND status IN ('completed', 'pending')
	`, today).Scan(&metrics.TodayDeposits.Count, &metrics.TodayDeposits.Volume)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Today's withdrawals
	err = h.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE type = 'withdrawal'
		AND created_at >= $1
		AND status IN ('completed', 'pending')
	`, today).Scan(&metrics.TodayWithdrawals.Count, &metrics.TodayWithdrawals.Volume)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Pending transactions
	err = h.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE status = 'pending'
	`).Scan(&metrics.Pending.Count, &metrics.Pending.Volume)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Failed transactions (today)
	err = h.db.QueryRowContext(ctx, `
		SELECT COUNT(*), COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE status = 'failed'
		AND created_at >= $1
	`, today).Scan(&metrics.Failed.Count, &metrics.Failed.Volume)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetPlatformWallets returns platform wallet statuses with 24h stats
// GET /api/v1/admin/transfers/wallets
func (h *TransferMonitoringHandler) GetPlatformWallets(c *gin.Context) {
	ctx := c.Request.Context()
	yesterday := time.Now().Add(-24 * time.Hour)

	query := `
		SELECT 
			pa.id, pa.currency, pa.account_type, pa.name, pa.balance,
			pa.min_balance, pa.max_balance, pa.is_active,
			COALESCE(dep.count, 0) as deposits_24h,
			COALESCE(wit.count, 0) as withdrawals_24h
		FROM platform_accounts pa
		LEFT JOIN (
			SELECT 
				pam.account_id,
				COUNT(*) as count
			FROM platform_account_movements pam
			WHERE pam.movement_type = 'debit'
			AND pam.created_at >= $1
			GROUP BY pam.account_id
		) dep ON pa.id = dep.account_id
		LEFT JOIN (
			SELECT 
				pam.account_id,
				COUNT(*) as count
			FROM platform_account_movements pam
			WHERE pam.movement_type = 'credit'
			AND pam.created_at >= $1
			GROUP BY pam.account_id
		) wit ON pa.id = wit.account_id
		WHERE pa.account_type = 'operations'
		ORDER BY pa.currency
	`

	rows, err := h.db.QueryContext(ctx, query, yesterday)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var wallets []WalletStatus
	for rows.Next() {
		var w WalletStatus
		err := rows.Scan(
			&w.ID, &w.Currency, &w.AccountType, &w.Name, &w.Balance,
			&w.MinBalance, &w.MaxBalance, &w.IsActive,
			&w.Deposits24h, &w.Withdrawals24h,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		wallets = append(wallets, w)
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}

// GetRecentTransactions returns recent transactions with pagination
// GET /api/v1/admin/transfers/transactions
func (h *TransferMonitoringHandler) GetRecentTransactions(c *gin.Context) {
	ctx := c.Request.Context()

	txType := c.Query("type")   // deposit, withdrawal
	status := c.Query("status") // pending, completed, failed
	limit := c.DefaultQuery("limit", "50")

	query := `
		SELECT 
			t.id, t.user_id, t.wallet_id, t.type, t.amount, t.currency,
			t.status, t.reference, t.provider, t.destination_account,
			t.created_at, t.updated_at
		FROM transactions t
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if txType != "" {
		query += ` AND t.type = $` + string(rune(argIndex+'0'))
		args = append(args, txType)
		argIndex++
	}

	if status != "" {
		query += ` AND t.status = $` + string(rune(argIndex+'0'))
		args = append(args, status)
		argIndex++
	}

	query += ` ORDER BY t.created_at DESC LIMIT $` + string(rune(argIndex+'0'))
	args = append(args, limit)

	rows, err := h.db.QueryContext(ctx, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(
			&t.ID, &t.UserID, &t.WalletID, &t.Type, &t.Amount, &t.Currency,
			&t.Status, &t.Reference, &t.Provider, &t.DestinationAccount,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		transactions = append(transactions, t)
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// GetAlerts returns critical alerts for platform wallets
// GET /api/v1/admin/transfers/alerts
func (h *TransferMonitoringHandler) GetAlerts(c *gin.Context) {
	ctx := c.Request.Context()

	query := `
		SELECT 
			pa.id, pa.currency, pa.name, pa.balance, pa.min_balance,
			CASE 
				WHEN pa.balance < pa.min_balance THEN 'critical'
				WHEN pa.balance < pa.min_balance * 2 THEN 'warning'
				ELSE 'normal'
			END as severity
		FROM platform_accounts pa
		WHERE pa.account_type = 'operations'
		AND pa.balance < pa.min_balance * 2
		ORDER BY 
			CASE 
				WHEN pa.balance < pa.min_balance THEN 1
				ELSE 2
			END,
			pa.balance ASC
	`

	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var alerts []WalletAlert
	for rows.Next() {
		var a WalletAlert
		err := rows.Scan(&a.ID, &a.Currency, &a.Name, &a.Balance, &a.MinBalance, &a.Severity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		alerts = append(alerts, a)
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

// Response types
type MetricData struct {
	Count  int     `json:"count"`
	Volume float64 `json:"volume"`
}

type WalletStatus struct {
	ID             string  `json:"id"`
	Currency       string  `json:"currency"`
	AccountType    string  `json:"account_type"`
	Name           string  `json:"name"`
	Balance        float64 `json:"balance"`
	MinBalance     float64 `json:"min_balance"`
	MaxBalance     float64 `json:"max_balance"`
	IsActive       bool    `json:"is_active"`
	Deposits24h    int     `json:"deposits_24h"`
	Withdrawals24h int     `json:"withdrawals_24h"`
}

type Transaction struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	WalletID           string    `json:"wallet_id"`
	Type               string    `json:"type"`
	Amount             float64   `json:"amount"`
	Currency           string    `json:"currency"`
	Status             string    `json:"status"`
	Reference          string    `json:"reference"`
	Provider           string    `json:"provider"`
	DestinationAccount *string   `json:"destination_account,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type WalletAlert struct {
	ID         string  `json:"id"`
	Currency   string  `json:"currency"`
	Name       string  `json:"name"`
	Balance    float64 `json:"balance"`
	MinBalance float64 `json:"min_balance"`
	Severity   string  `json:"severity"` // critical, warning, normal
}
