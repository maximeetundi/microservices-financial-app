package grpc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/repository"
)

// Server wraps a simple HTTP-based gRPC-like server
// Note: This is a lightweight implementation that avoids the grpc package
// due to Go version requirements. In production, upgrade Go and use real gRPC.
type Server struct {
	httpServer    *http.Server
	configService *ConfigService
	queryService  *QueryService
	repo          *repository.AdminRepository
}

// NewServer creates a new gRPC-like server (HTTP-based)
func NewServer(repo *repository.AdminRepository) *Server {
	configService := NewConfigService(repo)
	queryService := NewQueryService(repo)

	return &Server{
		configService: configService,
		queryService:  queryService,
		repo:          repo,
	}
}

// Start starts the server on the given port
func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	// Config service endpoints
	mux.HandleFunc("/grpc/config/get", s.handleGetConfig)
	mux.HandleFunc("/grpc/config/all", s.handleGetAllConfigs)
	mux.HandleFunc("/grpc/config/prefix", s.handleGetConfigsByPrefix)
	mux.HandleFunc("/grpc/feature", s.handleIsFeatureEnabled)
	mux.HandleFunc("/grpc/fee", s.handleGetFee)
	mux.HandleFunc("/grpc/limit", s.handleGetLimit)

	// Query service endpoints
	mux.HandleFunc("/grpc/dashboard", s.handleGetDashboardStats)
	mux.HandleFunc("/grpc/user/status", s.handleGetUserStatus)
	mux.HandleFunc("/grpc/platform/accounts", s.handleGetPlatformAccounts)
	mux.HandleFunc("/grpc/hot-wallets", s.handleGetHotWallets)

	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("[gRPC-HTTP] Admin service gRPC-like server starting on port %s", port)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[gRPC-HTTP] Server error: %v", err)
		}
	}()

	return nil
}

// Stop gracefully stops the server
func (s *Server) Stop() {
	log.Println("[gRPC-HTTP] Stopping server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.httpServer.Shutdown(ctx)
}

// HTTP handlers that proxy to service methods
func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	resp, err := s.configService.GetConfiguration(r.Context(), &GetConfigRequest{Key: key})
	s.writeResponse(w, resp, err)
}

func (s *Server) handleGetAllConfigs(w http.ResponseWriter, r *http.Request) {
	resp, err := s.configService.GetAllConfigurations(r.Context(), &GetAllConfigsRequest{})
	s.writeResponse(w, resp, err)
}

func (s *Server) handleGetConfigsByPrefix(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("prefix")
	resp, err := s.configService.GetConfigurationsByPrefix(r.Context(), &GetConfigsByPrefixRequest{Prefix: prefix})
	s.writeResponse(w, resp, err)
}

func (s *Server) handleIsFeatureEnabled(w http.ResponseWriter, r *http.Request) {
	feature := r.URL.Query().Get("feature")
	resp, err := s.configService.IsFeatureEnabled(r.Context(), &FeatureRequest{FeatureKey: feature})
	s.writeResponse(w, resp, err)
}

func (s *Server) handleGetFee(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	currency := r.URL.Query().Get("currency")
	resp, err := s.configService.GetFee(r.Context(), &GetFeeRequest{FeeKey: key, Currency: currency})
	s.writeResponse(w, resp, err)
}

func (s *Server) handleGetLimit(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	resp, err := s.configService.GetLimit(r.Context(), &GetLimitRequest{LimitKey: key})
	s.writeResponse(w, resp, err)
}

func (s *Server) handleGetDashboardStats(w http.ResponseWriter, r *http.Request) {
	resp, err := s.queryService.GetDashboardStats(r.Context(), &EmptyRequest{})
	s.writeResponse(w, resp, err)
}

func (s *Server) handleGetUserStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	resp, err := s.queryService.GetUserStatus(r.Context(), &GetUserStatusRequest{UserId: userID})
	s.writeResponse(w, resp, err)
}

func (s *Server) handleGetPlatformAccounts(w http.ResponseWriter, r *http.Request) {
	resp, err := s.queryService.GetPlatformAccounts(r.Context(), &EmptyRequest{})
	s.writeResponse(w, resp, err)
}

func (s *Server) handleGetHotWallets(w http.ResponseWriter, r *http.Request) {
	resp, err := s.queryService.GetHotWallets(r.Context(), &EmptyRequest{})
	s.writeResponse(w, resp, err)
}

func (s *Server) writeResponse(w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(data)
}

// ========== Message Types ==========

type GetConfigRequest struct {
	Key string `json:"key"`
}

type ConfigResponse struct {
	Key              string  `json:"key"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Type             string  `json:"type"`
	FixedAmount      float64 `json:"fixed_amount"`
	PercentageAmount float64 `json:"percentage_amount"`
	Currency         string  `json:"currency"`
	IsEnabled        bool    `json:"is_enabled"`
	UpdatedAt        string  `json:"updated_at"`
	UpdatedBy        string  `json:"updated_by"`
}

type GetAllConfigsRequest struct {
	Service string `json:"service"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

type ConfigListResponse struct {
	Configurations []*ConfigResponse `json:"configurations"`
	Total          int32             `json:"total"`
}

type GetConfigsByPrefixRequest struct {
	Prefix string `json:"prefix"`
}

type FeatureRequest struct {
	FeatureKey string `json:"feature_key"`
}

type FeatureResponse struct {
	IsEnabled bool `json:"is_enabled"`
}

type GetFeeRequest struct {
	FeeKey   string `json:"fee_key"`
	Currency string `json:"currency"`
}

type FeeResponse struct {
	FeeKey           string  `json:"fee_key"`
	Type             string  `json:"type"`
	FixedAmount      float64 `json:"fixed_amount"`
	PercentageAmount float64 `json:"percentage_amount"`
	Currency         string  `json:"currency"`
	IsEnabled        bool    `json:"is_enabled"`
}

type GetLimitRequest struct {
	LimitKey string `json:"limit_key"`
	KycLevel int32  `json:"kyc_level"`
}

type LimitResponse struct {
	LimitKey    string  `json:"limit_key"`
	LimitAmount float64 `json:"limit_amount"`
	Currency    string  `json:"currency"`
	TierName    string  `json:"tier_name"`
}

type EmptyRequest struct{}

type DashboardStatsResponse struct {
	TotalUsers        int64   `json:"total_users"`
	TotalWallets      int64   `json:"total_wallets"`
	TotalTransactions int64   `json:"total_transactions"`
	PendingKyc        int64   `json:"pending_kyc"`
	ActiveCards       int64   `json:"active_cards"`
	TotalVolume24h    float64 `json:"total_volume_24h"`
	Currency          string  `json:"currency"`
}

type GetUserStatusRequest struct {
	UserId string `json:"user_id"`
}

type UserStatusResponse struct {
	UserId        string `json:"user_id"`
	Email         string `json:"email"`
	Status        string `json:"status"`
	KycLevel      int32  `json:"kyc_level"`
	KycStatus     string `json:"kyc_status"`
	IsBlocked     bool   `json:"is_blocked"`
	BlockedReason string `json:"blocked_reason"`
	CreatedAt     string `json:"created_at"`
}

type KYCStatusResponse struct {
	UserId    string         `json:"user_id"`
	KycLevel  int32          `json:"kyc_level"`
	KycStatus string         `json:"kyc_status"`
	Documents []*KYCDocument `json:"documents"`
}

type KYCDocument struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	UploadedAt string `json:"uploaded_at"`
}

type GetWalletStatusRequest struct {
	WalletId string `json:"wallet_id"`
}

type WalletStatusResponse struct {
	WalletId         string  `json:"wallet_id"`
	UserId           string  `json:"user_id"`
	Currency         string  `json:"currency"`
	WalletType       string  `json:"wallet_type"`
	Status           string  `json:"status"`
	Balance          float64 `json:"balance"`
	AvailableBalance float64 `json:"available_balance"`
}

type PlatformAccountsResponse struct {
	Accounts []*PlatformAccount `json:"accounts"`
}

type PlatformAccount struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type HotWalletsResponse struct {
	Wallets []*HotWallet `json:"wallets"`
}

type HotWallet struct {
	Id       string  `json:"id"`
	Currency string  `json:"currency"`
	Address  string  `json:"address"`
	Balance  float64 `json:"balance"`
	Provider string  `json:"provider"`
}
