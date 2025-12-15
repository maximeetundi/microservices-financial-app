package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/repository"
	"github.com/streadway/amqp"
)

type WalletService struct {
	walletRepo      *repository.WalletRepository
	transactionRepo *repository.TransactionRepository
	cryptoService   *CryptoService
	balanceService  *BalanceService
	mqChannel       *amqp.Channel
}

func NewWalletService(
	walletRepo *repository.WalletRepository,
	transactionRepo *repository.TransactionRepository,
	cryptoService *CryptoService,
	balanceService *BalanceService,
	mqChannel *amqp.Channel,
) *WalletService {
	return &WalletService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		cryptoService:   cryptoService,
		balanceService:  balanceService,
		mqChannel:       mqChannel,
	}
}

func (s *WalletService) CreateWallet(userID string, req *models.CreateWalletRequest) (*models.Wallet, error) {
	// Check if user already has a wallet for this currency
	existingWallet, _ := s.walletRepo.GetByUserAndCurrency(userID, req.Currency)
	if existingWallet != nil {
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
	}

	// Generate crypto address if it's a crypto wallet
	if req.WalletType == "crypto" {
		cryptoWallet, err := s.cryptoService.GenerateWallet(req.Currency)
		if err != nil {
			return nil, fmt.Errorf("failed to generate crypto wallet: %w", err)
		}

		// Encrypt private key
		encryptedPrivateKey, err := s.cryptoService.EncryptPrivateKey(cryptoWallet.PrivateKey, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt private key: %w", err)
		}

		wallet.WalletAddress = &cryptoWallet.Address
		wallet.PrivateKeyEncrypted = &encryptedPrivateKey
	}

	err := s.walletRepo.Create(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Send wallet creation notification
	s.publishWalletEvent("wallet.created", wallet)

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

func (s *WalletService) GetUserWallets(userID string) ([]*models.Wallet, error) {
	wallets, err := s.walletRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Remove sensitive data
	for _, wallet := range wallets {
		wallet.PrivateKeyEncrypted = nil
	}

	return wallets, nil
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

	// Validate sufficient balance
	totalAmount := req.Amount
	estimatedFee, err := s.cryptoService.EstimateTransactionFee(wallet.Currency, req.Amount, "normal")
	if err == nil {
		totalAmount += estimatedFee.EstimatedFee
	}

	err = s.balanceService.ValidateSufficientBalance(walletID, totalAmount, false)
	if err != nil {
		return nil, err
	}

	// Validate destination address
	if !s.cryptoService.ValidateAddress(req.ToAddress, wallet.Currency) {
		return nil, fmt.Errorf("invalid destination address")
	}

	// Freeze the amount
	err = s.balanceService.FreezeAmount(walletID, totalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to freeze amount: %w", err)
	}

	// Create pending transaction
	transaction := &models.Transaction{
		FromWalletID:    &walletID,
		TransactionType: "send",
		Amount:          req.Amount,
		Fee:             estimatedFee.EstimatedFee,
		Currency:        wallet.Currency,
		Status:          "pending",
		Description:     req.Note,
	}

	// Create metadata
	metadata := map[string]interface{}{
		"to_address": req.ToAddress,
		"gas_price":  req.GasPrice,
	}
	metadataJSON, _ := json.Marshal(metadata)
	metadataStr := string(metadataJSON)
	transaction.Metadata = &metadataStr

	err = s.transactionRepo.Create(transaction)
	if err != nil {
		// Unfreeze amount on error
		s.balanceService.UnfreezeAmount(walletID, totalAmount)
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Create blockchain transaction
	txHash, err := s.cryptoService.CreateTransaction(wallet, req.ToAddress, req.Amount, req.GasPrice)
	if err != nil {
		// Update transaction status to failed
		s.transactionRepo.UpdateStatus(transaction.ID, "failed", nil)
		s.balanceService.UnfreezeAmount(walletID, totalAmount)
		return nil, fmt.Errorf("failed to create blockchain transaction: %w", err)
	}

	// Update transaction with blockchain hash
	transaction.BlockchainTxHash = &txHash
	s.transactionRepo.UpdateStatus(transaction.ID, "pending", &txHash)

	// Publish transaction event
	s.publishTransactionEvent("transaction.created", transaction)

	return transaction, nil
}

func (s *WalletService) ProcessCryptoDeposit(address, txHash, currency string, amount float64) error {
	// Find wallet by address
	// This would require a new repository method to find wallet by address
	// For now, we'll simulate this

	// Get wallet by address (mock implementation)
	walletID := "mock-wallet-id" // In reality, query by address

	// Create deposit transaction
	transaction := &models.Transaction{
		ToWalletID:       &walletID,
		TransactionType:  "deposit",
		Amount:           amount,
		Fee:              0,
		Currency:         strings.ToUpper(currency),
		Status:           "pending",
		BlockchainTxHash: &txHash,
	}

	err := s.transactionRepo.Create(transaction)
	if err != nil {
		return fmt.Errorf("failed to create deposit transaction: %w", err)
	}

	// Check confirmations and update balance when confirmed
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

	// Check balance
	if fromWallet.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance: %.2f available, %.2f required", fromWallet.Balance, req.Amount)
	}

	// Deduct from source
	if err := s.walletRepo.UpdateBalanceWithTransaction(req.FromWalletID, req.Amount, "send"); err != nil {
		return nil, fmt.Errorf("failed to deduct from source: %w", err)
	}

	// Add to destination
	if err := s.walletRepo.UpdateBalanceWithTransaction(req.ToWalletID, req.Amount, "receive"); err != nil {
		// Rollback source
		s.walletRepo.UpdateBalanceWithTransaction(req.FromWalletID, req.Amount, "receive")
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
		Fee:             0,
		Currency:        fromWallet.Currency,
		Status:          "completed",
		Description:     &description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

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
	if s.mqChannel == nil {
		return
	}

	event := map[string]interface{}{
		"type":      eventType,
		"wallet_id": wallet.ID,
		"user_id":   wallet.UserID,
		"currency":  wallet.Currency,
		"timestamp": time.Now(),
	}

	eventJSON, _ := json.Marshal(event)

	s.mqChannel.Publish(
		"wallet.events", // exchange
		eventType,       // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)
}

func (s *WalletService) publishTransactionEvent(eventType string, transaction *models.Transaction) {
	if s.mqChannel == nil {
		return
	}

	event := map[string]interface{}{
		"type":           eventType,
		"transaction_id": transaction.ID,
		"amount":         transaction.Amount,
		"currency":       transaction.Currency,
		"status":         transaction.Status,
		"timestamp":      time.Now(),
	}

	eventJSON, _ := json.Marshal(event)

	s.mqChannel.Publish(
		"transaction.events", // exchange
		eventType,            // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)
}