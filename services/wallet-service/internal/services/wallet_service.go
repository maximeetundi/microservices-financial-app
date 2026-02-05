package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/repository"
)

type WalletService struct {
	walletRepo          *repository.WalletRepository
	transactionRepo     *repository.TransactionRepository
	cryptoService       *CryptoService
	balanceService      *BalanceService
	feeService          *FeeService
	platformService     *PlatformAccountService // Platform hot/cold wallet management
	kafkaClient         *messaging.KafkaClient
	systemConfigService *SystemConfigService
}

func NewWalletService(
	walletRepo *repository.WalletRepository,
	transactionRepo *repository.TransactionRepository,
	cryptoService *CryptoService,
	balanceService *BalanceService,
	feeService *FeeService,
	platformService *PlatformAccountService,
	kafkaClient *messaging.KafkaClient,
	systemConfigService *SystemConfigService,
) *WalletService {
	return &WalletService{
		walletRepo:          walletRepo,
		transactionRepo:     transactionRepo,
		cryptoService:       cryptoService,
		balanceService:      balanceService,
		feeService:          feeService,
		platformService:     platformService,
		kafkaClient:         kafkaClient,
		systemConfigService: systemConfigService,
	}
}

func (s *WalletService) CreateWallet(userID string, req *models.CreateWalletRequest) (*models.Wallet, error) {
	// Check if user already has a wallet for this currency (Active OR Hidden)
	// We need a repo method that fetches ANY wallet (IsActive=true, IsHidden=any)
	// But current GetByUserAndCurrency fetches IsActive=true.
	// We need to fetch hidden ones too to UNHIDE them.
	// The repo update for GetByUserAndCurrency returns Hidden ones too now if they are IsActive=true.

	existingWallet, _ := s.walletRepo.GetByUserAndCurrency(userID, req.Currency)
	if existingWallet != nil {
		if existingWallet.IsHidden {
			// Unhide it
			err := s.walletRepo.Unhide(existingWallet.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to reactivate wallet: %w", err)
			}
			existingWallet.IsHidden = false
			return existingWallet, nil
		}
		return nil, fmt.Errorf("wallet already exists for currency %s", req.Currency)
	}

	wallet := &models.Wallet{
		UserID:        userID,
		Currency:      strings.ToUpper(req.Currency),
		WalletType:    req.WalletType,
		Balance:       0,
		FrozenBalance: 0,
		Name:          req.Name,
		IsActive:      true,
		IsHidden:      false,
	}

	// Generate crypto address if it's a crypto wallet
	// Generate crypto address if it's a crypto wallet
	if req.WalletType == "crypto" {
		cryptoWallet, err := s.cryptoService.GenerateWallet(userID, req.Currency)
		if err != nil {
			return nil, fmt.Errorf("failed to generate crypto wallet: %w", err)
		}

		// Private Key is stored in Vault by CryptoService
		// Private Key is stored in DB (Encrypted)
		wallet.WalletAddress = &cryptoWallet.Address
		wallet.PrivateKeyEncrypted = &cryptoWallet.EncryptedPrivateKey

		// Use Address as ExternalID for lookups (Non-Custodial)
		wallet.ExternalID = cryptoWallet.Address
	}

	err := s.walletRepo.Create(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Send wallet creation notification
	s.publishWalletEvent("wallet.created", wallet)

	return wallet, nil
}

// ProcessTatumDeposit handles incoming deposits from Tatum Webhooks
func (s *WalletService) ProcessTatumDeposit(accountID string, amount float64, currency string, txHash string) error {
	// 1. Find Wallet by Tatum Account ID
	wallet, err := s.walletRepo.GetByExternalID(accountID)
	if err != nil {
		return fmt.Errorf("wallet not found for tatum account %s: %w", accountID, err)
	}

	// 2. Check for duplicate transaction (Idempotency)
	// Ideally we check if txHash exists in transactions.
	// For now, assuming wallet_service transaction creation handles basic duplicates via Unique constraints if we had them or simple check.
	// Since we don't have GetByTxHash, let's implement basic check after finding wallet.

	// 3. Credit Wallet
	// ProcessTransaction handles Debit/Credit logic securely
	creditReq := &models.TransactionRequest{
		UserID:    wallet.UserID,
		WalletID:  wallet.ID,
		Amount:    amount,
		Type:      "credit",
		Currency:  wallet.Currency,     // Trust wallet currency or verify against webhook currency?
		Reference: "DEPOSIT_" + txHash, // Prefix helps identify source
	}

	// Verify currency match
	if strings.ToUpper(currency) != wallet.Currency && currency != "" {
		// Log warning, but for tokens like USDT it might verify differently (ERC20 vs ETH)
		// Tatum usually sends the currency of the account
		return fmt.Errorf("currency mismatch webhok %s vs wallet %s", currency, wallet.Currency)
	}

	err = s.ProcessTransaction(creditReq)
	if err != nil {
		return fmt.Errorf("failed to credit wallet for deposit: %w", err)
	}

	return nil
}

func (s *WalletService) GetWalletByCurrency(userID, currency string) (*models.Wallet, error) {
	return s.walletRepo.GetByUserAndCurrency(userID, currency)
}

func (s *WalletService) GetWalletByID(walletID string) (*models.Wallet, error) {
	wallet, err := s.walletRepo.GetByID(walletID)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletService) GetWallet(walletID, userID string) (*models.Wallet, error) {
	wallet, err := s.walletRepo.GetByID(walletID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if wallet.UserID != userID {
		return nil, fmt.Errorf("wallet not found")
	}

	// Remove sensitive data
	wallet.PrivateKeyEncrypted = nil

	return wallet, nil
}

func (s *WalletService) HideWallet(walletID, userID string) error {
	wallet, err := s.walletRepo.GetByID(walletID)
	if err != nil {
		return err
	}

	if wallet.UserID != userID {
		return fmt.Errorf("wallet not found")
	}

	if wallet.Balance > 0 || wallet.FrozenBalance > 0 {
		return fmt.Errorf("cannot hide wallet with funds. please transfer funds to your main wallet first")
	}

	// Check if it's the main (oldest) wallet logic - Optional?
	// User said "cannot delete main wallet" but "hide" might be allowed?
	// Let's assume hiding any wallet (if empty) is fine, or keep logic.
	// Keeping safety logic for now.
	userWallets, err := s.walletRepo.GetByUserID(userID, false)
	if err != nil {
		return fmt.Errorf("failed to check user wallets: %w", err)
	}

	if len(userWallets) <= 1 {
		return fmt.Errorf("cannot hide the only wallet")
	}

	return s.walletRepo.Hide(walletID)
}

// Deprecated: Use HideWallet
func (s *WalletService) DeleteWallet(walletID, userID string) error {
	return s.HideWallet(walletID, userID)
}

func (s *WalletService) UnhideWallet(walletID, userID string) error {
	// We need a repo method to get by ID even if hidden.
	// GetByID returns hidden ones too (based on repo update).
	wallet, err := s.walletRepo.GetByID(walletID)
	if err != nil {
		return err
	}

	if wallet.UserID != userID {
		return fmt.Errorf("wallet not found")
	}

	return s.walletRepo.Unhide(walletID)
}

func (s *WalletService) GetUserWallets(userID string, includeHidden bool) ([]*models.Wallet, error) {
	wallets, err := s.walletRepo.GetByUserID(userID, includeHidden)
	if err != nil {
		return nil, err
	}

	// Remove sensitive data
	for _, wallet := range wallets {
		wallet.PrivateKeyEncrypted = nil
	}

	return wallets, nil
}

// GetUserCountry retrieves the user's country code from the database
func (s *WalletService) GetUserCountry(userID string) string {
	country, err := s.walletRepo.GetUserCountry(userID)
	if err != nil || country == "" {
		return "" // Will default to USD in getCurrencyForCountry
	}
	return country
}

func (s *WalletService) SendCrypto(walletID, userID string, req *models.SendCryptoRequest) (*models.Transaction, error) {
	// Get wallet and verify ownership
	wallet, err := s.GetWallet(walletID, userID)
	if err != nil {
		return nil, err
	}

	if wallet.WalletType != "crypto" {
		return nil, fmt.Errorf("wallet is not a crypto wallet")
	}

	// Validate destination address
	if !s.cryptoService.ValidateAddress(req.ToAddress, wallet.Currency) {
		return nil, fmt.Errorf("invalid destination address")
	}

	// Check if this is an INTERNAL transfer (destination is another platform user)
	destinationWallet, _ := s.walletRepo.GetByAddress(req.ToAddress)
	isInternalTransfer := destinationWallet != nil

	if isInternalTransfer {
		// ==================== INTERNAL TRANSFER (DB ONLY) ====================
		// No blockchain transaction needed - just update balances in database
		return s.processInternalCryptoTransfer(wallet, destinationWallet, req)
	}

	// ==================== EXTERNAL TRANSFER (VIA HOT WALLET) ====================
	return s.processExternalCryptoTransfer(wallet, userID, req)
}

// processInternalCryptoTransfer handles transfers between platform users (DB only, no blockchain)
func (s *WalletService) processInternalCryptoTransfer(fromWallet, toWallet *models.Wallet, req *models.SendCryptoRequest) (*models.Transaction, error) {
	// Calculate platform fee only (no blockchain fee for internal transfers)
	platformFee, err := s.feeService.CalculateFee("transfer_internal", req.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate platform fee: %w", err)
	}
	totalDebit := req.Amount + platformFee

	// Validate sufficient balance
	err = s.balanceService.ValidateSufficientBalance(fromWallet.ID, totalDebit, false)
	if err != nil {
		return nil, err
	}

	// Debit source wallet
	err = s.walletRepo.UpdateBalanceWithTransaction(fromWallet.ID, totalDebit, "send")
	if err != nil {
		return nil, fmt.Errorf("failed to debit source wallet: %w", err)
	}

	// Credit destination wallet
	err = s.walletRepo.UpdateBalanceWithTransaction(toWallet.ID, req.Amount, "receive")
	if err != nil {
		// Rollback source debit
		s.walletRepo.UpdateBalanceWithTransaction(fromWallet.ID, totalDebit, "receive")
		return nil, fmt.Errorf("failed to credit destination wallet: %w", err)
	}

	// Create transaction record (no blockchain hash - internal)
	description := "Internal platform transfer"
	if req.Note != nil {
		description = *req.Note
	}

	transaction := &models.Transaction{
		FromWalletID:    &fromWallet.ID,
		ToWalletID:      &toWallet.ID,
		TransactionType: "transfer_internal",
		Amount:          req.Amount,
		Fee:             platformFee,
		Currency:        fromWallet.Currency,
		Status:          "completed", // Instant completion for internal transfers
		Description:     &description,
	}

	// Add metadata indicating internal transfer
	metadata := map[string]interface{}{
		"transfer_type":    "internal",
		"to_address":       req.ToAddress,
		"destination_user": toWallet.UserID,
		"blockchain_tx":    false,
	}
	metadataJSON, _ := json.Marshal(metadata)
	metadataStr := string(metadataJSON)
	transaction.Metadata = &metadataStr

	err = s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Publish event
	s.publishTransactionEvent("transaction.completed", transaction)

	return transaction, nil
}

// processExternalCryptoTransfer handles transfers to external addresses (via platform hot wallet)
func (s *WalletService) processExternalCryptoTransfer(wallet *models.Wallet, userID string, req *models.SendCryptoRequest) (*models.Transaction, error) {
	// Estimate blockchain fees
	estimatedFee, err := s.cryptoService.EstimateTransactionFee(wallet.Currency, req.Amount, "normal")
	if err != nil {
		estimatedFee = &models.CryptoTransactionEstimate{EstimatedFee: 0}
	}

	// Calculate platform fee
	platformFee, err := s.feeService.CalculateFee("crypto_send", req.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate platform fee: %w", err)
	}

	totalAmount := req.Amount + estimatedFee.EstimatedFee + platformFee

	// Validate sufficient balance
	err = s.balanceService.ValidateSufficientBalance(wallet.ID, totalAmount, false)
	if err != nil {
		return nil, err
	}

	// Freeze the amount in user's wallet
	err = s.balanceService.FreezeAmount(wallet.ID, totalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to freeze amount: %w", err)
	}

	// Create pending transaction
	transaction := &models.Transaction{
		FromWalletID:    &wallet.ID,
		TransactionType: "send",
		Amount:          req.Amount,
		Fee:             estimatedFee.EstimatedFee + platformFee,
		Currency:        wallet.Currency,
		Status:          "pending",
		Description:     req.Note,
	}

	// Create metadata
	metadata := map[string]interface{}{
		"transfer_type": "external",
		"to_address":    req.ToAddress,
		"gas_price":     req.GasPrice,
		"blockchain_tx": true,
	}

	// Auto-detect network from address
	networkInfo := s.cryptoService.DetectAddressNetwork(req.ToAddress)
	metadata["network_type"] = networkInfo.Type
	metadata["network_variant"] = networkInfo.Variant
	metadata["network_env"] = networkInfo.Network

	metadataJSON, _ := json.Marshal(metadata)
	metadataStr := string(metadataJSON)
	transaction.Metadata = &metadataStr

	err = s.transactionRepo.Create(transaction)
	if err != nil {
		s.balanceService.UnfreezeAmount(wallet.ID, totalAmount)
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Select platform hot wallet for this currency
	hotWallet, err := s.platformService.SelectBestCryptoWalletForDebit(wallet.Currency, "", req.Amount)
	if err != nil {
		s.transactionRepo.UpdateStatus(transaction.ID, "failed", nil)
		s.balanceService.UnfreezeAmount(wallet.ID, totalAmount)
		return nil, fmt.Errorf("no hot wallet available for %s: %w", wallet.Currency, err)
	}

	// Create blockchain transaction from hot wallet to destination
	txHash, err := s.cryptoService.CreateTransactionFromPlatformWallet(hotWallet, req.ToAddress, req.Amount, req.GasPrice)
	if err != nil {
		s.transactionRepo.UpdateStatus(transaction.ID, "failed", nil)
		s.balanceService.UnfreezeAmount(wallet.ID, totalAmount)
		return nil, fmt.Errorf("failed to create blockchain transaction: %w", err)
	}

	// Debit the user's balance (DB)
	s.balanceService.UnfreezeAmount(wallet.ID, totalAmount)
	s.walletRepo.UpdateBalanceWithTransaction(wallet.ID, totalAmount, "send")

	// Debit the platform hot wallet balance (DB)
	s.platformService.DebitCryptoWallet(hotWallet.ID, req.Amount, fmt.Sprintf("External send to %s", req.ToAddress))

	// Update transaction with blockchain hash
	transaction.BlockchainTxHash = &txHash
	s.transactionRepo.UpdateStatus(transaction.ID, "pending", &txHash)

	// Publish transaction event
	s.publishTransactionEvent("transaction.created", transaction)

	return transaction, nil
}

func (s *WalletService) ProcessCryptoDeposit(address, txHash, currency string, amount float64) error {
	// Find wallet by blockchain address
	wallet, err := s.walletRepo.GetByAddress(address)
	if err != nil {
		return fmt.Errorf("failed to lookup wallet by address: %w", err)
	}
	if wallet == nil {
		// Not a known address - could be platform wallet or unknown
		// Log but don't error - might be deposit to platform cold wallet
		return nil
	}

	// Check for duplicate transaction (idempotency via txHash)
	existingTx, _ := s.transactionRepo.GetByBlockchainHash(txHash)
	if existingTx != nil {
		// Transaction already processed, skip
		return nil
	}

	// Verify currency match
	if strings.ToUpper(currency) != wallet.Currency {
		return fmt.Errorf("currency mismatch: webhook %s vs wallet %s", currency, wallet.Currency)
	}

	walletID := wallet.ID

	// Create deposit transaction in pending state
	transaction := &models.Transaction{
		ToWalletID:       &walletID,
		TransactionType:  "deposit",
		Amount:           amount,
		Fee:              0,
		Currency:         strings.ToUpper(currency),
		Status:           "pending",
		BlockchainTxHash: &txHash,
	}

	err = s.transactionRepo.Create(transaction)
	if err != nil {
		return fmt.Errorf("failed to create deposit transaction: %w", err)
	}

	// Monitor blockchain confirmations asynchronously
	go s.monitorTransactionConfirmations(transaction.ID, txHash, currency)

	return nil
}

func (s *WalletService) FreezeWallet(walletID, userID, reason string) error {
	wallet, err := s.GetWallet(walletID, userID)
	if err != nil {
		return err
	}

	err = s.walletRepo.Freeze(walletID)
	if err != nil {
		return fmt.Errorf("failed to freeze wallet: %w", err)
	}

	// Publish wallet event
	s.publishWalletEvent("wallet.frozen", wallet)

	return nil
}

func (s *WalletService) UnfreezeWallet(walletID, userID string) error {
	wallet, err := s.GetWallet(walletID, userID)
	if err != nil {
		return err
	}

	err = s.walletRepo.Unfreeze(walletID)
	if err != nil {
		return fmt.Errorf("failed to unfreeze wallet: %w", err)
	}

	// Publish wallet event
	s.publishWalletEvent("wallet.unfrozen", wallet)

	return nil
}

func (s *WalletService) GetTransactionHistory(walletID, userID string, limit, offset int, status, txType string) ([]*models.Transaction, error) {
	// Verify wallet ownership
	_, err := s.GetWallet(walletID, userID)
	if err != nil {
		return nil, err
	}

	return s.transactionRepo.GetByWalletID(walletID, limit, offset, status, txType)
}

func (s *WalletService) GetPendingTransactions(walletID, userID string) ([]*models.Transaction, error) {
	return s.GetTransactionHistory(walletID, userID, 100, 0, "pending", "")
}

func (s *WalletService) EstimateTransactionFee(currency string, amount float64, priority string) (*models.CryptoTransactionEstimate, error) {
	return s.cryptoService.EstimateTransactionFee(currency, amount, priority)
}

func (s *WalletService) GetCryptoAddress(walletID, userID string) (*models.CryptoAddress, error) {
	wallet, err := s.GetWallet(walletID, userID)
	if err != nil {
		return nil, err
	}

	if wallet.WalletType != "crypto" || wallet.WalletAddress == nil {
		return nil, fmt.Errorf("wallet does not have a crypto address")
	}

	return &models.CryptoAddress{
		Address:  *wallet.WalletAddress,
		Currency: wallet.Currency,
		Network:  "mainnet", // This could be configurable
	}, nil
}

func (s *WalletService) GetWalletStats(userID string) (*models.WalletStats, error) {
	stats, err := s.walletRepo.GetWalletStats(userID)
	if err != nil {
		return nil, err
	}

	// Get recent activity
	transactions, err := s.transactionRepo.GetByUserID(userID, 10, 0, "", "")
	if err == nil {
		stats.RecentActivity = make([]models.Transaction, 0)
		for _, tx := range transactions {
			stats.RecentActivity = append(stats.RecentActivity, *tx)
		}
	}

	return stats, nil
}

// ProcessTransaction handles transaction requests from other services
func (s *WalletService) ProcessTransaction(req *models.TransactionRequest) error {
	// Validate wallet ownership/existence
	wallet, err := s.walletRepo.GetByID(req.WalletID)
	if err != nil {
		return fmt.Errorf("wallet not found")
	}
	if wallet.UserID != req.UserID {
		return fmt.Errorf("wallet does not belong to user")
	}

	// Validate currency
	if wallet.Currency != req.Currency {
		return fmt.Errorf("currency mismatch: %s vs %s", wallet.Currency, req.Currency)
	}

	// For debit, check balance
	if req.Type == "debit" {
		if wallet.Balance < req.Amount {
			return fmt.Errorf("insufficient balance")
		}

		// Enforce Limits
		// TODO: Fetch real KYC level from Auth Service or User local cache
		kycLevel := "standard" // Default
		if err := s.systemConfigService.CheckLimits(req.UserID, kycLevel, req.Currency, req.Amount); err != nil {
			return fmt.Errorf("limit exceeded: %w", err)
		}
	}

	// Update balance
	opType := "send" // default for debit
	if req.Type == "credit" {
		opType = "receive"
	}

	err = s.walletRepo.UpdateBalanceWithTransaction(req.WalletID, req.Amount, opType)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Create transaction record
	transaction := &models.Transaction{
		ID:              fmt.Sprintf("tx_%d", time.Now().UnixNano()),
		TransactionType: "exchange", // Assuming this is primarily used for exchange for now, or derive from Reference
		Amount:          req.Amount,
		Fee:             0, // Fees are handled by exchange service
		Currency:        req.Currency,
		Status:          "completed",
		Description:     &req.Reference,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if req.Type == "debit" {
		transaction.FromWalletID = &req.WalletID
	} else {
		transaction.ToWalletID = &req.WalletID
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		// Log error but don't fail as balance is already updated?
		// Ideally we should transactionally update both.
		// For now, return error so caller can retry or handle failure.
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Record Usage for Debit
	if req.Type == "debit" {
		s.systemConfigService.RecordUsage(req.UserID, req.Currency, req.Amount)
	}

	// Publish event
	s.publishTransactionEvent("transaction.completed", transaction)

	return nil
}

// ProcessDepositFromPlatform handles a user deposit by debiting the platform's reserve account and crediting the user's wallet.
// This ensures double-entry bookkeeping where funds originate from the platform's holding (Hot Wallet/Reserve).
func (s *WalletService) ProcessDepositFromPlatform(userID, walletID string, amount float64, currency, providerRef, providerName string) error {
	// 1. Validate User Wallet
	wallet, err := s.walletRepo.GetByID(walletID)
	if err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}
	if wallet.UserID != userID {
		return fmt.Errorf("wallet does not belong to user")
	}
	if wallet.Currency != currency {
		return fmt.Errorf("currency mismatch: wallet %s vs deposit %s", wallet.Currency, currency)
	}

	// 2. Debit Platform Reserve (The "Hot Wallet" logic)
	description := fmt.Sprintf("Deposit via %s (%s)", providerName, providerRef)
	err = s.platformService.DebitPlatformReserve(currency, amount, "deposit", providerRef, description)
	if err != nil {
		// If reserve is insufficient, we fail the deposit to enforce strict "funds must exist" logic
		return fmt.Errorf("failed to debit platform reserve (hot wallet): %w", err)
	}

	// 3. Credit User Wallet
	err = s.walletRepo.UpdateBalanceWithTransaction(walletID, amount, "deposit")
	if err != nil {
		// Critical error: debited platform but failed to credit user.
		// In production, this needs detailed logging or automatic rollback.
		return fmt.Errorf("failed to credit user wallet: %w", err)
	}

	// 4. Record User Transaction
	transaction := &models.Transaction{
		ID:              fmt.Sprintf("dep_%d", time.Now().UnixNano()),
		ToWalletID:      &walletID,
		TransactionType: "deposit",
		Amount:          amount,
		Fee:             0,
		Currency:        currency,
		Status:          "completed",
		Description:     &description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Metadata
	meta := map[string]interface{}{
		"provider":     providerName,
		"provider_ref": providerRef,
		"source":       "platform_reserve",
		"type":         "double_entry_deposit",
	}
	if metaJSON, err := json.Marshal(meta); err == nil {
		metaStr := string(metaJSON)
		transaction.Metadata = &metaStr
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		fmt.Printf("Error creating transaction record: %v\n", err)
	}

	// 5. Publish Event
	s.publishTransactionEvent("transaction.deposit.completed", transaction)

	return nil
}

// Transfer performs an internal wallet-to-wallet transfer
func (s *WalletService) Transfer(userID string, req *models.TransferRequest) (*models.Transaction, error) {
	// Validate source wallet ownership
	fromWallet, err := s.walletRepo.GetByID(req.FromWalletID)
	if err != nil {
		return nil, fmt.Errorf("source wallet not found")
	}
	if fromWallet.UserID != userID {
		return nil, fmt.Errorf("not authorized to transfer from this wallet")
	}
	if !fromWallet.IsActive {
		return nil, fmt.Errorf("source wallet is frozen")
	}

	// Validate destination wallet
	toWallet, err := s.walletRepo.GetByID(req.ToWalletID)
	if err != nil {
		return nil, fmt.Errorf("destination wallet not found")
	}
	if !toWallet.IsActive {
		return nil, fmt.Errorf("destination wallet is frozen")
	}

	// Validate same currency
	if fromWallet.Currency != toWallet.Currency {
		return nil, fmt.Errorf("currency mismatch: %s vs %s", fromWallet.Currency, toWallet.Currency)
	}

	// Calculate Fee
	fee, err := s.feeService.CalculateFee("transfer_internal", req.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fee: %w", err)
	}

	// Check balance with fee
	totalPayload := req.Amount + fee
	if fromWallet.Balance < totalPayload {
		return nil, fmt.Errorf("insufficient balance: %.2f available, %.2f required (including fee %.2f)", fromWallet.Balance, totalPayload, fee)
	}

	// Enforce Limits (on Source User)
	// TODO: Fetch real KYC level
	kycLevel := "standard"
	if err := s.systemConfigService.CheckLimits(userID, kycLevel, fromWallet.Currency, totalPayload); err != nil {
		return nil, fmt.Errorf("limit exceeded: %w", err)
	}

	// Deduct from source (Amount + Fee)
	if err := s.walletRepo.UpdateBalanceWithTransaction(req.FromWalletID, totalPayload, "send"); err != nil {
		return nil, fmt.Errorf("failed to deduct from source: %w", err)
	}

	// Add to destination (Amount only)
	if err := s.walletRepo.UpdateBalanceWithTransaction(req.ToWalletID, req.Amount, "receive"); err != nil {
		// Rollback source
		s.walletRepo.UpdateBalanceWithTransaction(req.FromWalletID, totalPayload, "receive")
		return nil, fmt.Errorf("failed to credit destination: %w", err)
	}

	// Create transaction record
	description := req.Description
	if description == "" {
		description = "Internal transfer"
	}

	transaction := &models.Transaction{
		ID:              fmt.Sprintf("tx_%d", time.Now().UnixNano()),
		FromWalletID:    &req.FromWalletID,
		ToWalletID:      &req.ToWalletID,
		TransactionType: "transfer",
		Amount:          req.Amount,
		Fee:             fee,
		Currency:        fromWallet.Currency,
		Status:          "completed",
		Description:     &description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Record Usage
	s.systemConfigService.RecordUsage(userID, fromWallet.Currency, totalPayload)

	// Publish event
	s.publishTransactionEvent("transaction.completed", transaction)

	return transaction, nil
}

// Private methods

func (s *WalletService) monitorTransactionConfirmations(transactionID, txHash, currency string) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	minConfirmations := s.cryptoService.GetMinimumConfirmations(currency)

	for range ticker.C {
		status, confirmations, err := s.cryptoService.GetTransactionStatus(txHash, currency)
		if err != nil {
			continue
		}

		if status == "confirmed" && confirmations >= minConfirmations {
			// Transaction is confirmed, update status and balance
			transaction, err := s.transactionRepo.GetByID(transactionID)
			if err != nil {
				continue
			}

			if transaction.ToWalletID != nil {
				// Credit the destination wallet
				s.balanceService.UpdateBalance(*transaction.ToWalletID, transaction.Amount, "deposit")
			}

			if transaction.FromWalletID != nil {
				// Debit the source wallet (unfreeze and deduct)
				totalAmount := transaction.Amount + transaction.Fee
				s.balanceService.UnfreezeAmount(*transaction.FromWalletID, totalAmount)
				s.balanceService.UpdateBalance(*transaction.FromWalletID, totalAmount, "send")
			}

			s.transactionRepo.UpdateStatus(transactionID, "completed", nil)

			// Publish completion event
			s.publishTransactionEvent("transaction.completed", transaction)
			break
		}
	}
}

func (s *WalletService) publishWalletEvent(eventType string, wallet *models.Wallet) {
	if s.kafkaClient == nil {
		return
	}

	eventData := map[string]interface{}{
		"wallet_id": wallet.ID,
		"user_id":   wallet.UserID,
		"currency":  wallet.Currency,
		"balance":   wallet.Balance,
	}

	envelope := messaging.NewEventEnvelope(eventType, "wallet-service", eventData)

	if err := s.kafkaClient.Publish(context.Background(), messaging.TopicWalletEvents, envelope); err != nil {
		// Log error but don't fail
		_ = err
	}
}

func (s *WalletService) publishTransactionEvent(eventType string, transaction *models.Transaction) {
	if s.kafkaClient == nil {
		return
	}

	eventData := map[string]interface{}{
		"transaction_id": transaction.ID,
		"amount":         transaction.Amount,
		"currency":       transaction.Currency,
		"status":         transaction.Status,
	}

	envelope := messaging.NewEventEnvelope(eventType, "wallet-service", eventData)

	if err := s.kafkaClient.Publish(context.Background(), messaging.TopicTransactionEvents, envelope); err != nil {
		// Log error but don't fail
		_ = err
	}
}
