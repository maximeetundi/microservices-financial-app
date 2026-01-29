package services

import (
	"fmt"
	"log"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/repository"
)

type PlatformAccountService struct {
	repo   *repository.PlatformAccountRepository
	crypto *CryptoService
}

func NewPlatformAccountService(repo *repository.PlatformAccountRepository, crypto *CryptoService) *PlatformAccountService {
	return &PlatformAccountService{
		repo:   repo,
		crypto: crypto,
	}
}

// Initialize seeds default platform accounts if they don't exist
func (s *PlatformAccountService) Initialize() error {
	log.Println("[Platform] Initializing platform accounts (Dual Wallet Architecture)...")

	// --- 1. Fiat Accounts (Storage & Operational) ---
	currencies := []string{"EUR", "USD", "GBP"}

	for _, currency := range currencies {
		// A. Storage Account (Reserve) - 1 Billion
		existsReserve, _ := s.repo.GetAccountByCurrencyType(currency, models.AccountTypeReserve)
		if existsReserve == nil {
			log.Printf("[Platform] Creating Reserve (Storage) Account for %s", currency)
			acct, err := s.CreateAccount(&models.CreatePlatformAccountRequest{
				Currency:    currency,
				AccountType: models.AccountTypeReserve,
				Name:        fmt.Sprintf("Reserve %s Storage", currency),
				MinBalance:  1000000.0,
				MaxBalance:  0, // Unlimited
				Priority:    100,
				Description: "Cold storage / Reserve funds",
			})
			if err == nil && acct != nil {
				// Seed with 1 Billion
				s.AdminCreditAccount(acct.ID, 1000000000, "Initial Reserve Seeding (1B)", "system_init", "genesis_reserve")
			}
		}

		// B. Operational Account (Hot) - 500 Million
		existsOps, _ := s.repo.GetAccountByCurrencyType(currency, models.AccountTypeOperations)
		if existsOps == nil {
			log.Printf("[Platform] Creating Operational Account for %s", currency)
			acct, err := s.CreateAccount(&models.CreatePlatformAccountRequest{
				Currency:    currency,
				AccountType: models.AccountTypeOperations,
				Name:        fmt.Sprintf("Operational %s Wallet", currency),
				MinBalance:  10000.0,
				MaxBalance:  0,
				Priority:    90,
				Description: "Active transactional funds",
			})
			if err == nil && acct != nil {
				// Seed with 500 Million
				s.AdminCreditAccount(acct.ID, 500000000, "Initial Ops Seeding (500M)", "system_init", "genesis_ops")
			}
		}

		// C. Fees Account (Accumulator)
		existsFee, _ := s.repo.GetAccountByCurrencyType(currency, models.AccountTypeFees)
		if existsFee == nil {
			log.Printf("[Platform] Creating Fee Account for %s", currency)
			s.CreateAccount(&models.CreatePlatformAccountRequest{
				Currency:    currency,
				AccountType: models.AccountTypeFees,
				Name:        fmt.Sprintf("Fees %s Accumulator", currency),
				Priority:    100,
			})
		}
	}

	// --- 2. Crypto Wallets (Cold & Hot) ---
	// Comprehensive list including Mainnet and Testnet variants
	cryptoCurrencies := []string{
		"BTC", "ETH", "SOL", "USDT", "BNB", "TRX", "MATIC",
		"TON", "XRP", "LTC", "DOGE", "BCH", "USDC", // Major L1s & Stablecoins
		"AVAX", "LINK", "UNI", "SHIB", "DAI", // Popular Alts/DeFi
		"BTC_TEST", "ETH_TEST", "SOL_TEST", "BNB_TEST", "MATIC_TEST", // Explicit Testnets
	}

	for _, currency := range cryptoCurrencies {
		// A. Cold Wallet (Storage)
		// Check for existing Cold wallet
		coldWallets, _ := s.repo.GetCryptoWalletsByCurrencyType(currency, models.WalletTypeCold)
		if len(coldWallets) == 0 {
			log.Printf("[Platform] Generating Cold Storage Wallet for %s...", currency)

			// User ID "platform_cold" for Vault segregation
			generated, err := s.crypto.GenerateWallet("platform_cold", currency)
			if err == nil {
				_, err = s.CreateCryptoWallet(&models.CreatePlatformCryptoWalletRequest{
					Currency:   currency,
					Network:    generated.Network,
					WalletType: models.WalletTypeCold, // "Stocker"
					Address:    generated.Address,
					Label:      fmt.Sprintf("%s Cold Storage", currency),
					MinBalance: 0,
					MaxBalance: 0,
					Priority:   100, // High priority for storage visibility
				})
				if err != nil {
					log.Printf("Error registering cold wallet for %s: %v", currency, err)
				}
			} else {
				log.Printf("Error generating cold keys for %s: %v", currency, err)
			}
		}

		// B. Hot Wallet (Operational)
		// Check for existing Hot wallet
		hotWallets, _ := s.repo.GetCryptoWalletsByCurrencyType(currency, models.WalletTypeHot)
		if len(hotWallets) == 0 {
			log.Printf("[Platform] Generating Hot Operational Wallet for %s...", currency)

			// User ID "platform_hot" for Vault segregation
			generated, err := s.crypto.GenerateWallet("platform_hot", currency)
			if err == nil {
				_, err = s.CreateCryptoWallet(&models.CreatePlatformCryptoWalletRequest{
					Currency:   currency,
					Network:    generated.Network,
					WalletType: models.WalletTypeHot, // "Utiliser"
					Address:    generated.Address,
					Label:      fmt.Sprintf("%s Ops Wallet 1", currency),
					MinBalance: 0,
					MaxBalance: 0,
					Priority:   90,
				})
				if err != nil {
					log.Printf("Error registering hot wallet for %s: %v", currency, err)
				}
			} else {
				log.Printf("Error generating hot keys for %s: %v", currency, err)
			}
		}
	}

	return nil
}

// ==================== Platform Fiat Accounts ====================

func (s *PlatformAccountService) GetAllAccounts() ([]models.PlatformAccount, error) {
	return s.repo.GetAllAccounts()
}

func (s *PlatformAccountService) GetAccountByID(id string) (*models.PlatformAccount, error) {
	return s.repo.GetAccountByID(id)
}

func (s *PlatformAccountService) GetReserveAccount(currency string) (*models.PlatformAccount, error) {
	return s.repo.GetAccountByCurrencyType(currency, models.AccountTypeReserve)
}

func (s *PlatformAccountService) GetFeeAccount(currency string) (*models.PlatformAccount, error) {
	return s.repo.GetAccountByCurrencyType(currency, models.AccountTypeFees)
}

func (s *PlatformAccountService) CreateAccount(req *models.CreatePlatformAccountRequest) (*models.PlatformAccount, error) {
	priority := req.Priority
	if priority == 0 {
		priority = 50 // Default priority
	}
	account := &models.PlatformAccount{
		Currency:    req.Currency,
		AccountType: req.AccountType,
		Name:        req.Name,
		Description: req.Description,
		Balance:     0,
		MinBalance:  req.MinBalance,
		MaxBalance:  req.MaxBalance,
		Priority:    priority,
	}
	err := s.repo.CreateAccount(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// AdminCreditAccount allows admin to manually credit a platform fiat account
func (s *PlatformAccountService) AdminCreditAccount(accountID string, amount float64, description, adminID, reference string) error {
	account, err := s.repo.GetAccountByID(accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return fmt.Errorf("account not found")
	}

	// Credit the account
	err = s.repo.CreditAccount(accountID, amount)
	if err != nil {
		return fmt.Errorf("failed to credit account: %w", err)
	}

	// Record the transaction
	tx := &models.PlatformTransaction{
		DebitAccountID:    "external",
		DebitAccountType:  models.AccTypeExternal,
		CreditAccountID:   accountID,
		CreditAccountType: models.AccTypePlatformFiat,
		Amount:            amount,
		Currency:          account.Currency,
		OperationType:     models.OpTypeAdminCredit,
		ReferenceType:     "admin_operation",
		ReferenceID:       reference,
		Description:       description,
		PerformedBy:       adminID,
	}
	if err := s.repo.CreateTransaction(tx); err != nil {
		log.Printf("Warning: Failed to record admin credit transaction: %v", err)
	}

	log.Printf("[Platform] Admin %s credited %f %s to account %s (ref: %s)", adminID, amount, account.Currency, accountID, reference)
	return nil
}

// AdminDebitAccount allows admin to manually debit a platform fiat account
func (s *PlatformAccountService) AdminDebitAccount(accountID string, amount float64, description, adminID, reference string) error {
	account, err := s.repo.GetAccountByID(accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	if account == nil {
		return fmt.Errorf("account not found")
	}
	if account.Balance < amount {
		return fmt.Errorf("insufficient balance: available %f, requested %f", account.Balance, amount)
	}

	// Debit the account
	err = s.repo.DebitAccount(accountID, amount)
	if err != nil {
		return fmt.Errorf("failed to debit account: %w", err)
	}

	// Record the transaction
	tx := &models.PlatformTransaction{
		DebitAccountID:    accountID,
		DebitAccountType:  models.AccTypePlatformFiat,
		CreditAccountID:   "external",
		CreditAccountType: models.AccTypeExternal,
		Amount:            amount,
		Currency:          account.Currency,
		OperationType:     models.OpTypeAdminDebit,
		ReferenceType:     "admin_operation",
		ReferenceID:       reference,
		Description:       description,
		PerformedBy:       adminID,
	}
	if err := s.repo.CreateTransaction(tx); err != nil {
		log.Printf("Warning: Failed to record admin debit transaction: %v", err)
	}

	log.Printf("[Platform] Admin %s debited %f %s from account %s (ref: %s)", adminID, amount, account.Currency, accountID, reference)
	return nil
}

// ==================== Crypto Wallets ====================

func (s *PlatformAccountService) GetAllCryptoWallets() ([]models.PlatformCryptoWallet, error) {
	return s.repo.GetAllCryptoWallets()
}

func (s *PlatformAccountService) GetCryptoWalletByID(id string) (*models.PlatformCryptoWallet, error) {
	return s.repo.GetCryptoWalletByID(id)
}

func (s *PlatformAccountService) GetCryptoWalletsByCurrency(currency string) ([]models.PlatformCryptoWallet, error) {
	return s.repo.GetCryptoWalletsByCurrency(currency)
}

func (s *PlatformAccountService) CreateCryptoWallet(req *models.CreatePlatformCryptoWalletRequest) (*models.PlatformCryptoWallet, error) {
	priority := req.Priority
	if priority == 0 {
		priority = 50 // Default priority
	}
	wallet := &models.PlatformCryptoWallet{
		Currency:   req.Currency,
		Network:    req.Network,
		WalletType: req.WalletType,
		Address:    req.Address,
		Label:      req.Label,
		Balance:    0,
		MinBalance: req.MinBalance,
		MaxBalance: req.MaxBalance,
		Priority:   priority,
	}
	err := s.repo.CreateCryptoWallet(wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *PlatformAccountService) UpdateCryptoWalletBalance(walletID string, balance float64) error {
	return s.repo.UpdateCryptoWalletBalance(walletID, balance)
}

// ==================== Transaction Ledger ====================

func (s *PlatformAccountService) GetTransactions(limit, offset int) ([]models.PlatformTransaction, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	return s.repo.GetTransactions(limit, offset)
}

func (s *PlatformAccountService) GetTransactionsByReference(refType, refID string) ([]models.PlatformTransaction, error) {
	return s.repo.GetTransactionsByReference(refType, refID)
}

// RecordTransaction records a double-entry transaction
func (s *PlatformAccountService) RecordTransaction(tx *models.PlatformTransaction) error {
	return s.repo.CreateTransaction(tx)
}

// ==================== Exchange Double-Entry Operations (with Intelligent Selection) ====================

// CreditPlatformReserve credits the platform reserve account (e.g., when user buys crypto)
// Uses intelligent selection to pick the best account based on priority and capacity
func (s *PlatformAccountService) CreditPlatformReserve(currency string, amount float64, referenceType, referenceID, description string) error {
	// Intelligent selection: pick best account for receiving funds
	account, err := s.repo.SelectAccountForCredit(currency, models.AccountTypeReserve, amount)
	if err != nil || account == nil {
		return fmt.Errorf("no platform reserve account available for currency %s: %v", currency, err)
	}

	log.Printf("[Platform] Crediting %.2f %s to account %s (%s)", amount, currency, account.ID, account.Name)

	err = s.repo.CreditAccount(account.ID, amount)
	if err != nil {
		return err
	}

	// Record transaction
	tx := &models.PlatformTransaction{
		DebitAccountID:    "user",
		DebitAccountType:  models.AccTypeUserWallet,
		CreditAccountID:   account.ID,
		CreditAccountType: models.AccTypePlatformFiat,
		Amount:            amount,
		Currency:          currency,
		OperationType:     models.OpTypeExchange,
		ReferenceType:     referenceType,
		ReferenceID:       referenceID,
		Description:       description,
	}
	return s.repo.CreateTransaction(tx)
}

// DebitPlatformReserve debits the platform reserve account (e.g., when user sells crypto)
// Uses intelligent selection to pick the best account with sufficient funds
func (s *PlatformAccountService) DebitPlatformReserve(currency string, amount float64, referenceType, referenceID, description string) error {
	// Intelligent selection: pick best account with sufficient balance
	account, err := s.repo.SelectAccountForDebit(currency, models.AccountTypeReserve, amount)
	if err != nil {
		return fmt.Errorf("cannot debit platform reserve: %v", err)
	}

	log.Printf("[Platform] Debiting %.2f %s from account %s (%s)", amount, currency, account.ID, account.Name)

	err = s.repo.DebitAccount(account.ID, amount)
	if err != nil {
		return err
	}

	// Record transaction
	tx := &models.PlatformTransaction{
		DebitAccountID:    account.ID,
		DebitAccountType:  models.AccTypePlatformFiat,
		CreditAccountID:   "user",
		CreditAccountType: models.AccTypeUserWallet,
		Amount:            amount,
		Currency:          currency,
		OperationType:     models.OpTypeExchange,
		ReferenceType:     referenceType,
		ReferenceID:       referenceID,
		Description:       description,
	}
	return s.repo.CreateTransaction(tx)
}

// CreditPlatformFees credits fees to the platform fees account
func (s *PlatformAccountService) CreditPlatformFees(currency string, amount float64, referenceType, referenceID, description string) error {
	if amount <= 0 {
		return nil // No fee to record
	}

	// Intelligent selection for fees account too
	account, err := s.repo.SelectAccountForCredit(currency, models.AccountTypeFees, amount)
	if err != nil || account == nil {
		return fmt.Errorf("platform fees account not found for currency %s", currency)
	}

	err = s.repo.CreditAccount(account.ID, amount)
	if err != nil {
		return err
	}

	// Record transaction
	tx := &models.PlatformTransaction{
		DebitAccountID:    "user",
		DebitAccountType:  models.AccTypeUserWallet,
		CreditAccountID:   account.ID,
		CreditAccountType: models.AccTypePlatformFiat,
		Amount:            amount,
		Currency:          currency,
		OperationType:     models.OpTypeFee,
		ReferenceType:     referenceType,
		ReferenceID:       referenceID,
		Description:       description,
	}
	return s.repo.CreateTransaction(tx)
}

// GetReconciliationReport returns balance totals for reconciliation
func (s *PlatformAccountService) GetReconciliationReport() (map[string]float64, error) {
	return s.repo.GetTotalBalanceByCurrency()
}

// ==================== Smart Selection Wrappers ====================

// SelectBestAccountForCredit exposes intelligent selection for external use
func (s *PlatformAccountService) SelectBestAccountForCredit(currency, accountType string, amount float64) (*models.PlatformAccount, error) {
	return s.repo.SelectAccountForCredit(currency, accountType, amount)
}

// SelectBestAccountForDebit exposes intelligent selection for external use
func (s *PlatformAccountService) SelectBestAccountForDebit(currency, accountType string, amount float64) (*models.PlatformAccount, error) {
	return s.repo.SelectAccountForDebit(currency, accountType, amount)
}

// SelectBestCryptoWalletForCredit exposes intelligent crypto wallet selection
func (s *PlatformAccountService) SelectBestCryptoWalletForCredit(currency, network string, amount float64) (*models.PlatformCryptoWallet, error) {
	return s.repo.SelectCryptoWalletForCredit(currency, network, amount)
}

// SelectBestCryptoWalletForDebit exposes intelligent crypto wallet selection
func (s *PlatformAccountService) SelectBestCryptoWalletForDebit(currency, network string, amount float64) (*models.PlatformCryptoWallet, error) {
	return s.repo.SelectCryptoWalletForDebit(currency, network, amount)
}
