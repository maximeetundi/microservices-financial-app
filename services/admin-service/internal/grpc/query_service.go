package grpc

import (
	"context"
	"fmt"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/repository"
)

// QueryService implements the gRPC AdminQueryService
type QueryService struct {
	repo *repository.AdminRepository
}

// NewQueryService creates a new QueryService
func NewQueryService(repo *repository.AdminRepository) *QueryService {
	return &QueryService{repo: repo}
}

// GetDashboardStats returns dashboard statistics
func (s *QueryService) GetDashboardStats(ctx context.Context, req *EmptyRequest) (*DashboardStatsResponse, error) {
	stats, err := s.repo.GetDashboardStats()
	if err != nil {
		return nil, err
	}

	return &DashboardStatsResponse{
		TotalUsers:        getInt64(stats, "total_users"),
		TotalWallets:      getInt64(stats, "total_wallets"),
		TotalTransactions: getInt64(stats, "total_transactions"),
		PendingKyc:        getInt64(stats, "pending_kyc"),
		ActiveCards:       getInt64(stats, "active_cards"),
		TotalVolume24h:    getFloat64(stats, "total_volume_24h"),
		Currency:          "EUR",
	}, nil
}

// GetUserStatus returns the status of a user
func (s *QueryService) GetUserStatus(ctx context.Context, req *GetUserStatusRequest) (*UserStatusResponse, error) {
	user, err := s.repo.GetUserByID(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", req.UserId)
	}

	isBlocked := false
	blockedReason := ""
	if blocked, ok := user["is_blocked"].(bool); ok {
		isBlocked = blocked
	}
	if reason, ok := user["blocked_reason"].(string); ok {
		blockedReason = reason
	}

	return &UserStatusResponse{
		UserId:        req.UserId,
		Email:         getString(user, "email"),
		Status:        getString(user, "status"),
		KycLevel:      int32(getInt64(user, "kyc_level")),
		KycStatus:     getString(user, "kyc_status"),
		IsBlocked:     isBlocked,
		BlockedReason: blockedReason,
		CreatedAt:     getString(user, "created_at"),
	}, nil
}

// GetUserKYCStatus returns KYC status for a user
func (s *QueryService) GetUserKYCStatus(ctx context.Context, req *GetUserStatusRequest) (*KYCStatusResponse, error) {
	user, err := s.repo.GetUserByID(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", req.UserId)
	}

	docs, _ := s.repo.GetUserKYCDocuments(req.UserId)

	documents := make([]*KYCDocument, len(docs))
	for i, doc := range docs {
		documents[i] = &KYCDocument{
			Id:         getString(doc, "id"),
			Type:       getString(doc, "document_type"),
			Status:     getString(doc, "status"),
			UploadedAt: getString(doc, "uploaded_at"),
		}
	}

	return &KYCStatusResponse{
		UserId:    req.UserId,
		KycLevel:  int32(getInt64(user, "kyc_level")),
		KycStatus: getString(user, "kyc_status"),
		Documents: documents,
	}, nil
}

// GetWalletStatus returns the status of a wallet
func (s *QueryService) GetWalletStatus(ctx context.Context, req *GetWalletStatusRequest) (*WalletStatusResponse, error) {
	wallet, err := s.repo.GetWalletByID(req.WalletId)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %s", req.WalletId)
	}

	return &WalletStatusResponse{
		WalletId:         req.WalletId,
		UserId:           getString(wallet, "user_id"),
		Currency:         getString(wallet, "currency"),
		WalletType:       getString(wallet, "wallet_type"),
		Status:           getString(wallet, "status"),
		Balance:          getFloat64(wallet, "balance"),
		AvailableBalance: getFloat64(wallet, "available_balance"),
	}, nil
}

// GetPlatformAccounts returns all platform accounts
func (s *QueryService) GetPlatformAccounts(ctx context.Context, req *EmptyRequest) (*PlatformAccountsResponse, error) {
	accounts, err := s.repo.GetPlatformAccounts()
	if err != nil {
		return nil, err
	}

	result := make([]*PlatformAccount, len(accounts))
	for i, acc := range accounts {
		result[i] = &PlatformAccount{
			Id:       getString(acc, "id"),
			Name:     getString(acc, "name"),
			Type:     getString(acc, "type"),
			Currency: getString(acc, "currency"),
			Balance:  getFloat64(acc, "balance"),
		}
	}

	return &PlatformAccountsResponse{Accounts: result}, nil
}

// GetHotWallets returns all hot wallets
func (s *QueryService) GetHotWallets(ctx context.Context, req *EmptyRequest) (*HotWalletsResponse, error) {
	wallets, err := s.repo.GetHotWallets()
	if err != nil {
		return nil, err
	}

	result := make([]*HotWallet, len(wallets))
	for i, w := range wallets {
		result[i] = &HotWallet{
			Id:       getString(w, "id"),
			Currency: getString(w, "currency"),
			Address:  getString(w, "address"),
			Balance:  getFloat64(w, "balance"),
			Provider: getString(w, "provider"),
		}
	}

	return &HotWalletsResponse{Wallets: result}, nil
}

// Helper functions to safely extract values from map
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case int64:
			return val
		case int:
			return int64(val)
		case float64:
			return int64(val)
		}
	}
	return 0
}

func getFloat64(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case int64:
			return float64(val)
		case int:
			return float64(val)
		}
	}
	return 0
}
