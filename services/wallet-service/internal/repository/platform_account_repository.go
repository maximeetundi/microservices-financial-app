package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/wallet-service/internal/models"
	"github.com/google/uuid"
)

type PlatformAccountRepository struct {
	db *sql.DB
}

func NewPlatformAccountRepository(db *sql.DB) *PlatformAccountRepository {
	return &PlatformAccountRepository{db: db}
}

// InitSchema creates the platform account tables and seeds default data
func (r *PlatformAccountRepository) InitSchema() error {
	// Platform Fiat Accounts (NO UNIQUE constraint - multiple accounts per currency allowed)
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS platform_accounts (
			id VARCHAR(36) PRIMARY KEY,
			currency VARCHAR(10) NOT NULL,
			account_type VARCHAR(50) NOT NULL,
			name VARCHAR(100) NOT NULL,
			balance DECIMAL(20, 8) DEFAULT 0,
			min_balance DECIMAL(20, 8) DEFAULT 0,
			max_balance DECIMAL(20, 8) DEFAULT 0,
			priority INT DEFAULT 50,
			description TEXT,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create platform_accounts table: %w", err)
	}

	// Add new columns if table already exists (migration)
	r.db.Exec(`ALTER TABLE platform_accounts ADD COLUMN IF NOT EXISTS min_balance DECIMAL(20, 8) DEFAULT 0`)
	r.db.Exec(`ALTER TABLE platform_accounts ADD COLUMN IF NOT EXISTS max_balance DECIMAL(20, 8) DEFAULT 0`)
	r.db.Exec(`ALTER TABLE platform_accounts ADD COLUMN IF NOT EXISTS priority INT DEFAULT 50`)
	// Remove UNIQUE constraint if exists (PostgreSQL specific)
	r.db.Exec(`ALTER TABLE platform_accounts DROP CONSTRAINT IF EXISTS platform_accounts_currency_account_type_key`)
	if err != nil {
		return fmt.Errorf("failed to create platform_accounts table: %w", err)
	}

	// Platform Crypto Wallets (multiple wallets per currency/network allowed)
	// SECURITY: Private keys are stored encrypted, never in plaintext
	_, err = r.db.Exec(`
		CREATE TABLE IF NOT EXISTS platform_crypto_wallets (
			id VARCHAR(36) PRIMARY KEY,
			currency VARCHAR(10) NOT NULL,
			network VARCHAR(50) NOT NULL,
			wallet_type VARCHAR(50) NOT NULL,
			address VARCHAR(100) NOT NULL,
			label VARCHAR(100),
			balance DECIMAL(30, 18) DEFAULT 0,
			min_balance DECIMAL(30, 18) DEFAULT 0,
			max_balance DECIMAL(30, 18) DEFAULT 0,
			priority INT DEFAULT 50,
			encrypted_private_key TEXT,
			key_hash VARCHAR(64),
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create platform_crypto_wallets table: %w", err)
	}

	// Add encryption columns if table already exists (migration)
	r.db.Exec(`ALTER TABLE platform_crypto_wallets ADD COLUMN IF NOT EXISTS encrypted_private_key TEXT`)
	r.db.Exec(`ALTER TABLE platform_crypto_wallets ADD COLUMN IF NOT EXISTS key_hash VARCHAR(64)`)

	// Create index
	r.db.Exec(`CREATE INDEX IF NOT EXISTS idx_pcw_currency_network ON platform_crypto_wallets(currency, network)`)

	// Platform Transactions (Ledger)
	_, err = r.db.Exec(`
		CREATE TABLE IF NOT EXISTS platform_transactions (
			id VARCHAR(36) PRIMARY KEY,
			transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			debit_account_id VARCHAR(36),
			debit_account_type VARCHAR(20),
			credit_account_id VARCHAR(36),
			credit_account_type VARCHAR(20),
			amount DECIMAL(30, 18) NOT NULL,
			currency VARCHAR(10) NOT NULL,
			operation_type VARCHAR(50) NOT NULL,
			reference_type VARCHAR(50),
			reference_id VARCHAR(100),
			description TEXT,
			performed_by VARCHAR(100),
			blockchain_tx_hash VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create platform_transactions table: %w", err)
	}

	// Create indexes for transactions
	r.db.Exec(`CREATE INDEX IF NOT EXISTS idx_pt_operation ON platform_transactions(operation_type)`)
	r.db.Exec(`CREATE INDEX IF NOT EXISTS idx_pt_date ON platform_transactions(transaction_date)`)
	r.db.Exec(`CREATE INDEX IF NOT EXISTS idx_pt_reference ON platform_transactions(reference_type, reference_id)`)

	return r.seedDefaultAccounts()
}

// seedDefaultAccounts creates the default platform accounts and crypto wallets
func (r *PlatformAccountRepository) seedDefaultAccounts() error {
	// Initial balance for fiat reserve accounts: 1 billion per currency
	const INITIAL_FIAT_RESERVE = 1000000000.0 // 1 milliard

	// Seed fiat accounts with initial balance for reserves
	type SeedAccount struct {
		models.PlatformAccount
		InitialBalance float64
	}

	defaultAccounts := []SeedAccount{
		// FCFA accounts - reserve gets 1 billion
		{PlatformAccount: models.PlatformAccount{Currency: "FCFA", AccountType: "reserve", Name: "Réserve FCFA Principal", Description: "Réserve principale en FCFA", Priority: 100}, InitialBalance: INITIAL_FIAT_RESERVE},
		{PlatformAccount: models.PlatformAccount{Currency: "FCFA", AccountType: "fees", Name: "Frais collectés FCFA", Description: "Frais de transaction collectés", Priority: 100}, InitialBalance: 0},
		{PlatformAccount: models.PlatformAccount{Currency: "FCFA", AccountType: "operations", Name: "Opérations FCFA", Description: "Compte opérationnel pour retraits/dépôts", Priority: 80}, InitialBalance: 100000000}, // 100M pour opérations
		// EUR accounts
		{PlatformAccount: models.PlatformAccount{Currency: "EUR", AccountType: "reserve", Name: "Réserve EUR", Description: "Réserve principale en EUR", Priority: 100}, InitialBalance: INITIAL_FIAT_RESERVE},
		{PlatformAccount: models.PlatformAccount{Currency: "EUR", AccountType: "fees", Name: "Frais collectés EUR", Description: "Frais de transaction en EUR", Priority: 100}, InitialBalance: 0},
		// USD accounts
		{PlatformAccount: models.PlatformAccount{Currency: "USD", AccountType: "reserve", Name: "Réserve USD", Description: "Réserve principale en USD", Priority: 100}, InitialBalance: INITIAL_FIAT_RESERVE},
		{PlatformAccount: models.PlatformAccount{Currency: "USD", AccountType: "fees", Name: "Frais collectés USD", Description: "Frais de transaction en USD", Priority: 100}, InitialBalance: 0},
		// XOF accounts
		{PlatformAccount: models.PlatformAccount{Currency: "XOF", AccountType: "reserve", Name: "Réserve XOF", Description: "Réserve principale en XOF", Priority: 100}, InitialBalance: INITIAL_FIAT_RESERVE},
	}

	for _, acc := range defaultAccounts {
		var exists bool
		err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM platform_accounts WHERE currency = $1 AND account_type = $2 AND name = $3)",
			acc.Currency, acc.AccountType, acc.Name).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			id := uuid.New().String()
			_, err := r.db.Exec(`
				INSERT INTO platform_accounts (id, currency, account_type, name, description, balance, priority, min_balance, max_balance, is_active)
				VALUES ($1, $2, $3, $4, $5, $6, $7, 0, 0, true)
			`, id, acc.Currency, acc.AccountType, acc.Name, acc.Description, acc.InitialBalance, acc.Priority)
			if err != nil {
				return fmt.Errorf("failed to seed platform account %s/%s: %w", acc.Currency, acc.AccountType, err)
			}
			log.Printf("[Platform] Seeded account: %s %s with balance %.2f", acc.Currency, acc.AccountType, acc.InitialBalance)
		}
	}

	// Seed default crypto wallets (EMPTY - no address, no balance - must be configured by admin)
	// These are just placeholder entries for the admin to configure
	defaultCryptoWallets := []models.PlatformCryptoWallet{
		// Bitcoin - empty, admin configures addresses
		{Currency: "BTC", Network: "bitcoin", WalletType: "hot", Label: "BTC Hot Wallet Principal", Priority: 100},
		{Currency: "BTC", Network: "bitcoin", WalletType: "cold", Label: "BTC Cold Storage", Priority: 50},
		// Ethereum
		{Currency: "ETH", Network: "ethereum", WalletType: "hot", Label: "ETH Hot Wallet Principal", Priority: 100},
		{Currency: "ETH", Network: "ethereum", WalletType: "cold", Label: "ETH Cold Storage", Priority: 50},
		// USDT on multiple networks
		{Currency: "USDT", Network: "ethereum", WalletType: "hot", Label: "USDT ERC20 Hot Wallet", Priority: 100},
		{Currency: "USDT", Network: "tron", WalletType: "hot", Label: "USDT TRC20 Hot Wallet", Priority: 100},
		{Currency: "USDT", Network: "bsc", WalletType: "hot", Label: "USDT BEP20 Hot Wallet", Priority: 80},
		// USDC
		{Currency: "USDC", Network: "ethereum", WalletType: "hot", Label: "USDC ERC20 Hot Wallet", Priority: 100},
		// BNB
		{Currency: "BNB", Network: "bsc", WalletType: "hot", Label: "BNB Hot Wallet", Priority: 100},
		// SOL
		{Currency: "SOL", Network: "solana", WalletType: "hot", Label: "SOL Hot Wallet", Priority: 100},
	}

	for _, wallet := range defaultCryptoWallets {
		var exists bool
		err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM platform_crypto_wallets WHERE currency = $1 AND network = $2 AND label = $3)",
			wallet.Currency, wallet.Network, wallet.Label).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			id := uuid.New().String()
			// Address is empty, balance is 0 - must be configured by admin with encrypted private key
			_, err := r.db.Exec(`
				INSERT INTO platform_crypto_wallets (id, currency, network, wallet_type, address, label, balance, priority, min_balance, max_balance, is_active)
				VALUES ($1, $2, $3, $4, '', $5, 0, $6, 0, 0, true)
			`, id, wallet.Currency, wallet.Network, wallet.WalletType, wallet.Label, wallet.Priority)
			if err != nil {
				return fmt.Errorf("failed to seed crypto wallet %s/%s: %w", wallet.Currency, wallet.Network, err)
			}
			log.Printf("[Platform] Seeded crypto wallet placeholder: %s on %s (NO ADDRESS - configure in admin)", wallet.Currency, wallet.Network)
		}
	}

	return nil
}

// ==================== Platform Accounts CRUD ====================

func (r *PlatformAccountRepository) GetAllAccounts() ([]models.PlatformAccount, error) {
	query := `SELECT id, currency, account_type, name, balance, min_balance, max_balance, priority, description, is_active, created_at, updated_at 
			  FROM platform_accounts ORDER BY currency, priority DESC, account_type`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.PlatformAccount
	for rows.Next() {
		var a models.PlatformAccount
		err := rows.Scan(&a.ID, &a.Currency, &a.AccountType, &a.Name, &a.Balance, &a.MinBalance, &a.MaxBalance, &a.Priority, &a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (r *PlatformAccountRepository) GetAccountByID(id string) (*models.PlatformAccount, error) {
	query := `SELECT id, currency, account_type, name, balance, min_balance, max_balance, priority, description, is_active, created_at, updated_at 
			  FROM platform_accounts WHERE id = $1`
	var a models.PlatformAccount
	err := r.db.QueryRow(query, id).Scan(&a.ID, &a.Currency, &a.AccountType, &a.Name, &a.Balance, &a.MinBalance, &a.MaxBalance, &a.Priority, &a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *PlatformAccountRepository) GetAccountByCurrencyType(currency, accountType string) (*models.PlatformAccount, error) {
	query := `SELECT id, currency, account_type, name, balance, description, is_active, created_at, updated_at 
			  FROM platform_accounts WHERE currency = $1 AND account_type = $2`
	var a models.PlatformAccount
	err := r.db.QueryRow(query, currency, accountType).Scan(&a.ID, &a.Currency, &a.AccountType, &a.Name, &a.Balance, &a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *PlatformAccountRepository) CreateAccount(account *models.PlatformAccount) error {
	account.ID = uuid.New().String()
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()
	account.IsActive = true
	if account.Priority == 0 {
		account.Priority = 50
	}

	_, err := r.db.Exec(`
		INSERT INTO platform_accounts (id, currency, account_type, name, balance, min_balance, max_balance, priority, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`, account.ID, account.Currency, account.AccountType, account.Name, account.Balance, account.MinBalance, account.MaxBalance, account.Priority, account.Description, account.IsActive, account.CreatedAt, account.UpdatedAt)
	return err
}

func (r *PlatformAccountRepository) UpdateBalance(id string, newBalance float64) error {
	_, err := r.db.Exec(`UPDATE platform_accounts SET balance = $1, updated_at = NOW() WHERE id = $2`, newBalance, id)
	return err
}

func (r *PlatformAccountRepository) CreditAccount(id string, amount float64) error {
	_, err := r.db.Exec(`UPDATE platform_accounts SET balance = balance + $1, updated_at = NOW() WHERE id = $2`, amount, id)
	return err
}

func (r *PlatformAccountRepository) DebitAccount(id string, amount float64) error {
	_, err := r.db.Exec(`UPDATE platform_accounts SET balance = balance - $1, updated_at = NOW() WHERE id = $2`, amount, id)
	return err
}

// ==================== Crypto Wallets CRUD ====================

func (r *PlatformAccountRepository) GetAllCryptoWallets() ([]models.PlatformCryptoWallet, error) {
	query := `SELECT id, currency, network, wallet_type, address, label, balance, min_balance, max_balance, priority, is_active, created_at, updated_at 
			  FROM platform_crypto_wallets ORDER BY currency, priority DESC, network`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []models.PlatformCryptoWallet
	for rows.Next() {
		var w models.PlatformCryptoWallet
		err := rows.Scan(&w.ID, &w.Currency, &w.Network, &w.WalletType, &w.Address, &w.Label, &w.Balance, &w.MinBalance, &w.MaxBalance, &w.Priority, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}

func (r *PlatformAccountRepository) GetCryptoWalletByID(id string) (*models.PlatformCryptoWallet, error) {
	query := `SELECT id, currency, network, wallet_type, address, label, balance, min_balance, max_balance, priority, is_active, created_at, updated_at 
			  FROM platform_crypto_wallets WHERE id = $1`
	var w models.PlatformCryptoWallet
	err := r.db.QueryRow(query, id).Scan(&w.ID, &w.Currency, &w.Network, &w.WalletType, &w.Address, &w.Label, &w.Balance, &w.MinBalance, &w.MaxBalance, &w.Priority, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *PlatformAccountRepository) GetCryptoWalletsByCurrency(currency string) ([]models.PlatformCryptoWallet, error) {
	query := `SELECT id, currency, network, wallet_type, address, label, balance, min_balance, max_balance, priority, is_active, created_at, updated_at 
			  FROM platform_crypto_wallets WHERE currency = $1 AND is_active = true ORDER BY priority DESC, network`
	rows, err := r.db.Query(query, currency)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []models.PlatformCryptoWallet
	for rows.Next() {
		var w models.PlatformCryptoWallet
		err := rows.Scan(&w.ID, &w.Currency, &w.Network, &w.WalletType, &w.Address, &w.Label, &w.Balance, &w.MinBalance, &w.MaxBalance, &w.Priority, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}

func (r *PlatformAccountRepository) CreateCryptoWallet(wallet *models.PlatformCryptoWallet) error {
	wallet.ID = uuid.New().String()
	wallet.CreatedAt = time.Now()
	wallet.UpdatedAt = time.Now()
	wallet.IsActive = true
	if wallet.Priority == 0 {
		wallet.Priority = 50
	}

	_, err := r.db.Exec(`
		INSERT INTO platform_crypto_wallets (id, currency, network, wallet_type, address, label, balance, min_balance, max_balance, priority, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, wallet.ID, wallet.Currency, wallet.Network, wallet.WalletType, wallet.Address, wallet.Label, wallet.Balance, wallet.MinBalance, wallet.MaxBalance, wallet.Priority, wallet.IsActive, wallet.CreatedAt, wallet.UpdatedAt)
	return err
}

func (r *PlatformAccountRepository) UpdateCryptoWalletBalance(id string, balance float64) error {
	_, err := r.db.Exec(`UPDATE platform_crypto_wallets SET balance = $1, updated_at = NOW() WHERE id = $2`, balance, id)
	return err
}

func (r *PlatformAccountRepository) CreditCryptoWallet(id string, amount float64) error {
	_, err := r.db.Exec(`UPDATE platform_crypto_wallets SET balance = balance + $1, updated_at = NOW() WHERE id = $2`, amount, id)
	return err
}

func (r *PlatformAccountRepository) DebitCryptoWallet(id string, amount float64) error {
	_, err := r.db.Exec(`UPDATE platform_crypto_wallets SET balance = balance - $1, updated_at = NOW() WHERE id = $2`, amount, id)
	return err
}

// ==================== Platform Transactions (Ledger) ====================

func (r *PlatformAccountRepository) CreateTransaction(tx *models.PlatformTransaction) error {
	tx.ID = uuid.New().String()
	tx.TransactionDate = time.Now()
	tx.CreatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO platform_transactions (id, transaction_date, debit_account_id, debit_account_type, 
			credit_account_id, credit_account_type, amount, currency, operation_type, reference_type, 
			reference_id, description, performed_by, blockchain_tx_hash, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`, tx.ID, tx.TransactionDate, tx.DebitAccountID, tx.DebitAccountType, tx.CreditAccountID,
		tx.CreditAccountType, tx.Amount, tx.Currency, tx.OperationType, tx.ReferenceType,
		tx.ReferenceID, tx.Description, tx.PerformedBy, tx.BlockchainTxHash, tx.CreatedAt)
	return err
}

func (r *PlatformAccountRepository) GetTransactions(limit, offset int) ([]models.PlatformTransaction, error) {
	query := `SELECT id, transaction_date, debit_account_id, debit_account_type, credit_account_id, 
			  credit_account_type, amount, currency, operation_type, reference_type, reference_id, 
			  description, performed_by, blockchain_tx_hash, created_at 
			  FROM platform_transactions ORDER BY transaction_date DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.PlatformTransaction
	for rows.Next() {
		var t models.PlatformTransaction
		err := rows.Scan(&t.ID, &t.TransactionDate, &t.DebitAccountID, &t.DebitAccountType,
			&t.CreditAccountID, &t.CreditAccountType, &t.Amount, &t.Currency, &t.OperationType,
			&t.ReferenceType, &t.ReferenceID, &t.Description, &t.PerformedBy, &t.BlockchainTxHash, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *PlatformAccountRepository) GetTransactionsByReference(refType, refID string) ([]models.PlatformTransaction, error) {
	query := `SELECT id, transaction_date, debit_account_id, debit_account_type, credit_account_id, 
			  credit_account_type, amount, currency, operation_type, reference_type, reference_id, 
			  description, performed_by, blockchain_tx_hash, created_at 
			  FROM platform_transactions WHERE reference_type = $1 AND reference_id = $2 ORDER BY created_at`
	rows, err := r.db.Query(query, refType, refID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.PlatformTransaction
	for rows.Next() {
		var t models.PlatformTransaction
		err := rows.Scan(&t.ID, &t.TransactionDate, &t.DebitAccountID, &t.DebitAccountType,
			&t.CreditAccountID, &t.CreditAccountType, &t.Amount, &t.Currency, &t.OperationType,
			&t.ReferenceType, &t.ReferenceID, &t.Description, &t.PerformedBy, &t.BlockchainTxHash, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

// GetTotalBalance returns total balance for all platform fiat accounts by currency
func (r *PlatformAccountRepository) GetTotalBalanceByCurrency() (map[string]float64, error) {
	query := `SELECT currency, SUM(balance) as total FROM platform_accounts GROUP BY currency`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var currency string
		var total float64
		if err := rows.Scan(&currency, &total); err != nil {
			return nil, err
		}
		result[currency] = total
	}
	return result, nil
}

// ==================== Intelligent Account Selection ====================

// SelectAccountForCredit selects the best account to receive funds (deposit)
// Priority: highest priority first, then accounts with most capacity (furthest from max_balance)
func (r *PlatformAccountRepository) SelectAccountForCredit(currency, accountType string, amount float64) (*models.PlatformAccount, error) {
	// Query: active accounts ordered by priority DESC, then by available capacity DESC
	query := `
		SELECT id, currency, account_type, name, balance, min_balance, max_balance, priority, description, is_active, created_at, updated_at
		FROM platform_accounts 
		WHERE currency = $1 AND account_type = $2 AND is_active = true
		  AND (max_balance = 0 OR balance + $3 <= max_balance)
		ORDER BY priority DESC, (max_balance - balance) DESC
		LIMIT 1
	`
	var a models.PlatformAccount
	err := r.db.QueryRow(query, currency, accountType, amount).Scan(
		&a.ID, &a.Currency, &a.AccountType, &a.Name, &a.Balance, &a.MinBalance, &a.MaxBalance, &a.Priority, &a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		// Fallback: get any active account
		return r.GetAnyActiveAccount(currency, accountType)
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// SelectAccountForDebit selects the best account to source funds (withdrawal)
// Priority: highest priority first, then accounts with highest available balance (above min_balance)
func (r *PlatformAccountRepository) SelectAccountForDebit(currency, accountType string, amount float64) (*models.PlatformAccount, error) {
	// Query: active accounts with sufficient balance (above min_balance), ordered by priority DESC, then by balance DESC
	query := `
		SELECT id, currency, account_type, name, balance, min_balance, max_balance, priority, description, is_active, created_at, updated_at
		FROM platform_accounts 
		WHERE currency = $1 AND account_type = $2 AND is_active = true
		  AND (balance - $3) >= min_balance
		ORDER BY priority DESC, balance DESC
		LIMIT 1
	`
	var a models.PlatformAccount
	err := r.db.QueryRow(query, currency, accountType, amount).Scan(
		&a.ID, &a.Currency, &a.AccountType, &a.Name, &a.Balance, &a.MinBalance, &a.MaxBalance, &a.Priority, &a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no account with sufficient funds: need %.2f %s", amount, currency)
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetAnyActiveAccount returns any active account for given currency/type
func (r *PlatformAccountRepository) GetAnyActiveAccount(currency, accountType string) (*models.PlatformAccount, error) {
	query := `
		SELECT id, currency, account_type, name, balance, min_balance, max_balance, priority, description, is_active, created_at, updated_at
		FROM platform_accounts 
		WHERE currency = $1 AND account_type = $2 AND is_active = true
		ORDER BY priority DESC
		LIMIT 1
	`
	var a models.PlatformAccount
	err := r.db.QueryRow(query, currency, accountType).Scan(
		&a.ID, &a.Currency, &a.AccountType, &a.Name, &a.Balance, &a.MinBalance, &a.MaxBalance, &a.Priority, &a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetAccountsByCurrencyType returns ALL accounts for a currency/type (for display)
func (r *PlatformAccountRepository) GetAccountsByCurrencyType(currency, accountType string) ([]models.PlatformAccount, error) {
	query := `
		SELECT id, currency, account_type, name, balance, min_balance, max_balance, priority, description, is_active, created_at, updated_at
		FROM platform_accounts 
		WHERE currency = $1 AND account_type = $2 AND is_active = true
		ORDER BY priority DESC
	`
	rows, err := r.db.Query(query, currency, accountType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.PlatformAccount
	for rows.Next() {
		var a models.PlatformAccount
		err := rows.Scan(&a.ID, &a.Currency, &a.AccountType, &a.Name, &a.Balance, &a.MinBalance, &a.MaxBalance, &a.Priority, &a.Description, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

// ==================== Intelligent Crypto Wallet Selection ====================

// SelectCryptoWalletForCredit selects best crypto wallet to receive funds
func (r *PlatformAccountRepository) SelectCryptoWalletForCredit(currency, network string, amount float64) (*models.PlatformCryptoWallet, error) {
	query := `
		SELECT id, currency, network, wallet_type, address, label, balance, min_balance, max_balance, priority, is_active, created_at, updated_at
		FROM platform_crypto_wallets 
		WHERE currency = $1 AND network = $2 AND is_active = true
		  AND (max_balance = 0 OR balance + $3 <= max_balance)
		ORDER BY priority DESC, (max_balance - balance) DESC
		LIMIT 1
	`
	var w models.PlatformCryptoWallet
	err := r.db.QueryRow(query, currency, network, amount).Scan(
		&w.ID, &w.Currency, &w.Network, &w.WalletType, &w.Address, &w.Label, &w.Balance, &w.MinBalance, &w.MaxBalance, &w.Priority, &w.IsActive, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no crypto wallet available for %s on %s", currency, network)
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}

// SelectCryptoWalletForDebit selects best crypto wallet to source funds
func (r *PlatformAccountRepository) SelectCryptoWalletForDebit(currency, network string, amount float64) (*models.PlatformCryptoWallet, error) {
	query := `
		SELECT id, currency, network, wallet_type, address, label, balance, min_balance, max_balance, priority, is_active, created_at, updated_at
		FROM platform_crypto_wallets 
		WHERE currency = $1 AND network = $2 AND is_active = true
		  AND (balance - $3) >= min_balance
		ORDER BY priority DESC, balance DESC
		LIMIT 1
	`
	var w models.PlatformCryptoWallet
	err := r.db.QueryRow(query, currency, network, amount).Scan(
		&w.ID, &w.Currency, &w.Network, &w.WalletType, &w.Address, &w.Label, &w.Balance, &w.MinBalance, &w.MaxBalance, &w.Priority, &w.IsActive, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no crypto wallet with sufficient funds: need %.8f %s on %s", amount, currency, network)
	}
	if err != nil {
		return nil, err
	}
	return &w, nil
}
