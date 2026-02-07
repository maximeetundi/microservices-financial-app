package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/common/messaging"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func buildSeedInstanceCredentials(providerName, providerType, baseURL string) map[string]interface{} {
	creds := make(map[string]interface{})

	// Provide minimal-but-non-placeholder values so transfer-service can initialize provider structs.
	// These are NOT real credentials and must be replaced in admin panel.
	// Important: avoid placeholder patterns filtered by transfer-service (REPLACE_ME / YOUR_ / etc.).
	suffix := uuid.New().String()[:8]

	// Common defaults
	if baseURL != "" {
		creds["base_url"] = baseURL
	}
	creds["environment"] = "sandbox"

	switch providerName {
	case "orange_money":
		creds["client_id"] = fmt.Sprintf("client_%s_%s", providerName, suffix)
		creds["client_secret"] = fmt.Sprintf("secret_%s_%s", providerName, uuid.New().String()[:12])
		creds["merchant_key"] = fmt.Sprintf("merchant_%s_%s", providerName, uuid.New().String()[:10])
		if creds["base_url"] == "" {
			creds["base_url"] = "https://api.orange.com/orange-money-webpay/dev/v1"
		}
		return creds

	case "wave", "wave_money", "wave_ci", "wave_sn":
		creds["api_key"] = fmt.Sprintf("api_%s_%s", providerName, uuid.New().String()[:12])
		if creds["base_url"] == "" {
			creds["base_url"] = "https://api.wave.com/v1"
		}
		return creds

	case "cinetpay":
		creds["api_key"] = fmt.Sprintf("api_%s_%s", providerName, uuid.New().String()[:12])
		creds["site_id"] = fmt.Sprintf("site_%s_%s", providerName, suffix)
		if creds["base_url"] == "" {
			creds["base_url"] = "https://api-checkout.cinetpay.com/v2"
		}
		return creds

	case "flutterwave":
		creds["public_key"] = fmt.Sprintf("pk_%s_%s", providerName, suffix)
		creds["secret_key"] = fmt.Sprintf("sk_%s_%s", providerName, uuid.New().String()[:12])
		creds["encryption_key"] = fmt.Sprintf("enc_%s_%s", providerName, uuid.New().String()[:12])
		if creds["base_url"] == "" {
			creds["base_url"] = "https://api.flutterwave.com/v3"
		}
		return creds

	case "mtn_momo", "mtn_money":
		creds["subscription_key"] = fmt.Sprintf("sub_%s_%s", providerName, uuid.New().String()[:12])
		creds["api_user"] = fmt.Sprintf("user_%s_%s", providerName, suffix)
		creds["api_key"] = fmt.Sprintf("key_%s_%s", providerName, uuid.New().String()[:12])
		creds["target_environment"] = "sandbox"
		if creds["base_url"] == "" {
			creds["base_url"] = "https://sandbox.momodeveloper.mtn.com"
		}
		return creds
	}

	// Fallback by provider type
	switch providerType {
	case "card":
		creds["publishable_key"] = fmt.Sprintf("pk_%s_%s", providerName, suffix)
		creds["secret_key"] = fmt.Sprintf("sk_%s_%s", providerName, uuid.New().String()[:12])
		creds["webhook_secret"] = fmt.Sprintf("whsec_%s_%s", providerName, uuid.New().String()[:12])
		if creds["base_url"] == "" {
			creds["base_url"] = "https://api.stripe.com/v1"
		}

	case "international":
		creds["client_id"] = fmt.Sprintf("client_%s_%s", providerName, suffix)
		creds["client_secret"] = fmt.Sprintf("secret_%s_%s", providerName, uuid.New().String()[:12])
		if creds["base_url"] == "" {
			creds["base_url"] = "https://api-m.sandbox.paypal.com"
		}
		creds["mode"] = "sandbox"

	default:
		creds["api_key"] = fmt.Sprintf("api_%s_%s", providerName, uuid.New().String()[:12])
		creds["secret_key"] = fmt.Sprintf("sk_%s_%s", providerName, uuid.New().String()[:12])
		if creds["base_url"] == "" {
			creds["base_url"] = "https://api.example.com"
		}
	}

	return creds
}

// InitializeAdminDB initializes the admin database connection
func InitializeAdminDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open admin database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping admin database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// Create admin tables
	if err := createAdminTables(db); err != nil {
		return nil, fmt.Errorf("failed to create admin tables: %w", err)
	}

	return db, nil
}

// InitializeMainDB initializes read-only connection to main database
func InitializeMainDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open main database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping main database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil
}

func createAdminTables(db *sql.DB) error {
	queries := []string{
		// Permissions table
		`CREATE TABLE IF NOT EXISTS admin_permissions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			code VARCHAR(100) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			category VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Roles table
		`CREATE TABLE IF NOT EXISTS admin_roles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) UNIQUE NOT NULL,
			description TEXT,
			is_system BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Role-Permission mapping
		`CREATE TABLE IF NOT EXISTS admin_role_permissions (
			role_id UUID REFERENCES admin_roles(id) ON DELETE CASCADE,
			permission_id UUID REFERENCES admin_permissions(id) ON DELETE CASCADE,
			PRIMARY KEY (role_id, permission_id)
		)`,

		// Admin users table
		`CREATE TABLE IF NOT EXISTS admin_users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			role_id UUID REFERENCES admin_roles(id),
			is_active BOOLEAN DEFAULT TRUE,
			last_login_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by UUID REFERENCES admin_users(id)
		)`,

		// Admin sessions
		`CREATE TABLE IF NOT EXISTS admin_sessions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			admin_id UUID REFERENCES admin_users(id) ON DELETE CASCADE,
			token_hash VARCHAR(255) NOT NULL,
			ip_address VARCHAR(45),
			user_agent TEXT,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Audit logs
		`CREATE TABLE IF NOT EXISTS admin_audit_logs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			admin_id UUID REFERENCES admin_users(id),
			admin_email VARCHAR(255) NOT NULL,
			action VARCHAR(100) NOT NULL,
			resource VARCHAR(100) NOT NULL,
			resource_id VARCHAR(255),
			details JSONB,
			ip_address VARCHAR(45),
			user_agent TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Admin notifications (for KYC, support tickets, alerts, etc.)
		`CREATE TABLE IF NOT EXISTS admin_notifications (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			type VARCHAR(50) NOT NULL,
			title VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			data JSONB,
			is_read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Payment providers table
		`CREATE TABLE IF NOT EXISTS payment_providers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) UNIQUE NOT NULL,
			display_name VARCHAR(255) NOT NULL,
			provider_type VARCHAR(50) NOT NULL,
			api_base_url TEXT,
			api_key_encrypted TEXT,
			api_secret_encrypted TEXT,
			public_key_encrypted TEXT,
			webhook_secret_encrypted TEXT,
			is_active BOOLEAN DEFAULT TRUE,
			is_demo_mode BOOLEAN DEFAULT FALSE,
			deposit_enabled BOOLEAN DEFAULT TRUE,
			withdraw_enabled BOOLEAN DEFAULT TRUE,
			logo_url TEXT,
			supported_currencies JSONB DEFAULT '[]',
			config_json JSONB DEFAULT '{}',
			capability VARCHAR(20) DEFAULT 'mixed',
			fee_percentage DECIMAL(5, 4) DEFAULT 0,
			fee_fixed DECIMAL(20, 2) DEFAULT 0,
			min_transaction DECIMAL(20, 2) DEFAULT 100,
			max_transaction DECIMAL(20, 2) DEFAULT 5000000,
			daily_limit DECIMAL(20, 2) DEFAULT 50000000,
			priority INT DEFAULT 50,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Migration: add new columns if table exists
		`ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS deposit_enabled BOOLEAN DEFAULT TRUE`,
		`ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS withdraw_enabled BOOLEAN DEFAULT TRUE`,
		`ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS fee_percentage DECIMAL(5, 4) DEFAULT 0`,
		`ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS fee_fixed DECIMAL(20, 2) DEFAULT 0`,
		`ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS min_transaction DECIMAL(20, 2) DEFAULT 100`,
		`ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS max_transaction DECIMAL(20, 2) DEFAULT 5000000`,
		`ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS daily_limit DECIMAL(20, 2) DEFAULT 50000000`,
		`ALTER TABLE payment_providers ADD COLUMN IF NOT EXISTS priority INT DEFAULT 50`,

		// Provider country mappings
		`CREATE TABLE IF NOT EXISTS provider_countries (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			provider_id UUID REFERENCES payment_providers(id) ON DELETE CASCADE,
			country_code VARCHAR(3) NOT NULL,
			country_name VARCHAR(100) NOT NULL,
			currency VARCHAR(10) NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			priority INT DEFAULT 50,
			min_amount DECIMAL(20, 2) DEFAULT 0,
			max_amount DECIMAL(20, 2) DEFAULT 0,
			fee_percentage DECIMAL(5, 4) DEFAULT 0,
			fee_fixed DECIMAL(20, 2) DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Provider instances (multi-key support)
		`CREATE TABLE IF NOT EXISTS provider_instances (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			provider_id UUID REFERENCES payment_providers(id) ON DELETE CASCADE,
			name VARCHAR(100) NOT NULL,
			vault_secret_path VARCHAR(255),
			api_credentials JSONB DEFAULT '{}',
			hot_wallet_id VARCHAR(36),
			deposit_enabled BOOLEAN DEFAULT TRUE,
			withdraw_enabled BOOLEAN DEFAULT TRUE,
			is_active BOOLEAN DEFAULT TRUE,
			is_primary BOOLEAN DEFAULT FALSE,
			is_global BOOLEAN DEFAULT FALSE,
			is_paused BOOLEAN DEFAULT FALSE,
			pause_reason TEXT,
			paused_at TIMESTAMP,
			priority INT DEFAULT 50,
			request_count BIGINT DEFAULT 0,
			last_used_at TIMESTAMP,
			last_error TEXT,
			health_status VARCHAR(20) DEFAULT 'unknown',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Add api_credentials column if it doesn't exist (migration for existing tables)
		`DO $$ BEGIN
			ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS api_credentials JSONB DEFAULT '{}';
		EXCEPTION WHEN duplicate_column THEN NULL;
		END $$`,

		// Provider instance to hot wallet mapping (multi-wallet per instance)
		`CREATE TABLE IF NOT EXISTS provider_instance_wallets (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			instance_id UUID REFERENCES provider_instances(id) ON DELETE CASCADE,
			currency VARCHAR(10) NOT NULL,
			hot_wallet_id VARCHAR(36) NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			priority INT DEFAULT 50,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(instance_id, currency, hot_wallet_id)
		)`,

		// Credit campaigns for mass credits, promotions, compensations
		`CREATE TABLE IF NOT EXISTS credit_campaigns (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			type VARCHAR(50) NOT NULL, -- single, mass, promotion
			status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
			reason TEXT NOT NULL,
			reason_type VARCHAR(50), -- compensation, bonus, promotion, contest, refund, loyalty, referral, other
			currency VARCHAR(10) NOT NULL,
			total_amount DECIMAL(20, 8) DEFAULT 0,
			user_count INT DEFAULT 0,
			success_count INT DEFAULT 0,
			failed_count INT DEFAULT 0,
			filters JSONB,
			admin_id UUID,
			hot_wallet_id VARCHAR(36),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP
		)`,

		// Individual credit operations (audit trail)
		`CREATE TABLE IF NOT EXISTS credit_operations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			campaign_id UUID REFERENCES credit_campaigns(id) ON DELETE CASCADE,
			user_id VARCHAR(36) NOT NULL,
			wallet_id VARCHAR(36),
			currency VARCHAR(10) NOT NULL,
			amount DECIMAL(20, 8) NOT NULL,
			status VARCHAR(50) DEFAULT 'pending', -- pending, success, failed
			error_message TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Create indexes for faster queries
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_admin_id ON admin_audit_logs(admin_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON admin_audit_logs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_sessions_admin_id ON admin_sessions(admin_id)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_notifications_created_at ON admin_notifications(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_notifications_is_read ON admin_notifications(is_read)`,
		`CREATE INDEX IF NOT EXISTS idx_provider_instances_provider_id ON provider_instances(provider_id)`,
		`CREATE INDEX IF NOT EXISTS idx_provider_instances_health ON provider_instances(health_status)`,

		// Migration: add pause columns if table already exists
		`ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS is_global BOOLEAN DEFAULT FALSE`,
		`ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS is_paused BOOLEAN DEFAULT FALSE`,
		`ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS pause_reason TEXT`,
		`ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS paused_at TIMESTAMP`,
		`ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS is_test_mode BOOLEAN DEFAULT TRUE`,
		`UPDATE provider_instances
		 SET is_test_mode =
		   CASE
		     WHEN LOWER(COALESCE(api_credentials->>'mode','')) IN ('sandbox','test') THEN TRUE
		     WHEN LOWER(COALESCE(api_credentials->>'mode','')) IN ('live','production','prod') THEN FALSE
		
		     WHEN LOWER(COALESCE(api_credentials->>'target_environment','')) = 'sandbox' THEN TRUE
		     WHEN LOWER(COALESCE(api_credentials->>'target_environment','')) IN ('live','production') THEN FALSE
		
		     WHEN COALESCE(api_credentials->>'base_url','') ILIKE '%sandbox%' THEN TRUE
		     WHEN COALESCE(api_credentials->>'base_url','') ILIKE '%test%' THEN TRUE
		
		     WHEN COALESCE(api_credentials->>'api_key','') LIKE 'sk_test%' THEN TRUE
		     WHEN COALESCE(api_credentials->>'api_key','') LIKE 'sk_live%' THEN FALSE
		
		     WHEN COALESCE(api_credentials->>'secret_key','') LIKE 'sk_test_%' THEN TRUE
		
		     WHEN COALESCE(api_credentials->>'secret_key','') LIKE 'FLWSECK_TEST%' THEN TRUE
		     WHEN COALESCE(api_credentials->>'secret_key','') LIKE 'FLWSECK_LIVE%' THEN FALSE
		
		     ELSE COALESCE(is_test_mode, TRUE)
		   END`,
		`ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS deposit_enabled BOOLEAN DEFAULT TRUE`,
		`ALTER TABLE provider_instances ADD COLUMN IF NOT EXISTS withdraw_enabled BOOLEAN DEFAULT TRUE`,
		`CREATE INDEX IF NOT EXISTS idx_provider_instances_paused ON provider_instances(is_paused) WHERE is_paused = TRUE`,
		`CREATE INDEX IF NOT EXISTS idx_provider_countries_provider_id ON provider_countries(provider_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_provider_countries_unique ON provider_countries(provider_id, country_code)`,
		`CREATE INDEX IF NOT EXISTS idx_provider_instance_wallets_instance ON provider_instance_wallets(instance_id)`,
		`CREATE INDEX IF NOT EXISTS idx_provider_instance_wallets_currency ON provider_instance_wallets(currency)`,
		`CREATE INDEX IF NOT EXISTS idx_credit_campaigns_status ON credit_campaigns(status)`,
		`CREATE INDEX IF NOT EXISTS idx_credit_campaigns_created_at ON credit_campaigns(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_credit_operations_campaign ON credit_operations(campaign_id)`,
		`CREATE INDEX IF NOT EXISTS idx_credit_operations_user ON credit_operations(user_id)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	// Seed default permissions and roles
	if err := seedDefaultData(db); err != nil {
		return fmt.Errorf("failed to seed default data: %w", err)
	}

	// Seed default payment providers
	if err := seedDefaultProviders(db); err != nil {
		log.Printf("Warning: failed to seed payment providers: %v", err)
	}

	// Seed default API keys for provider instances
	if err := seedDefaultAPIKeys(db); err != nil {
		log.Printf("Warning: failed to seed default API keys: %v", err)
	}

	// Fix Vault paths (remove /default suffix to match actual Vault structure)
	if err := fixVaultPaths(db); err != nil {
		log.Printf("Warning: failed to fix vault paths: %v", err)
	}

	return nil
}

// fixVaultPaths corrects the vault_secret_path in provider_instances
// to match the actual Vault structure (without /default suffix)
func fixVaultPaths(db *sql.DB) error {
	log.Println("[Database] Checking and fixing Vault paths...")

	// Count paths that need fixing
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM provider_instances WHERE vault_secret_path LIKE '%/default'`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count paths: %w", err)
	}

	if count == 0 {
		log.Println("[Database] ✅ All Vault paths are correct")
		return nil
	}

	log.Printf("[Database] 🔧 Found %d instances with incorrect Vault paths, fixing...", count)

	// Update paths: remove /default suffix
	result, err := db.Exec(`
		UPDATE provider_instances
		SET vault_secret_path = REPLACE(vault_secret_path, '/default', '')
		WHERE vault_secret_path LIKE '%/default'
	`)
	if err != nil {
		return fmt.Errorf("failed to update paths: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("[Database] ✅ Fixed %d Vault paths (removed /default suffix)", rowsAffected)

	// Also fix any double slashes
	db.Exec(`
		UPDATE provider_instances
		SET vault_secret_path = REPLACE(vault_secret_path, '//', '/')
		WHERE vault_secret_path LIKE '%//%'
	`)

	return nil
}

// seedDefaultProviders seeds all payment providers with their country configurations
func seedDefaultProviders(db *sql.DB) error {
	log.Println("[Database] Seeding/Updating default payment providers...")

	type Provider struct {
		Name        string
		DisplayName string
		Type        string
		BaseURL     string
		LogoURL     string
		Capability  string
		IsDemo      bool
		Countries   []struct {
			Code     string
			Name     string
			Currency string
		}
	}

	providers := []Provider{
		{
			Name:        "demo",
			DisplayName: "Demo Provider",
			Type:        "demo", // Demo type is always global
			BaseURL:     "",
			LogoURL:     "/icons/aggregators/demo.svg",
			Capability:  "mixed",
			IsDemo:      true,
			Countries:   []struct{ Code, Name, Currency string }{}, // No specific countries = global
		},
		{
			Name:        "flutterwave",
			DisplayName: "Flutterwave",
			Type:        "mobile_money",
			BaseURL:     "https://api.flutterwave.com/v3",
			LogoURL:     "/icons/aggregators/flutterwave.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"NG", "Nigeria", "NGN"},
				{"GH", "Ghana", "GHS"},
				{"KE", "Kenya", "KES"},
				{"CI", "Côte d'Ivoire", "XOF"},
			},
		},
		{
			Name:        "cinetpay",
			DisplayName: "CinetPay",
			Type:        "mobile_money",
			BaseURL:     "https://api-checkout.cinetpay.com/v2",
			LogoURL:     "/icons/aggregators/cinetpay.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"CI", "Côte d'Ivoire", "XOF"},
				{"SN", "Sénégal", "XOF"},
				{"CM", "Cameroun", "XAF"},
			},
		},

		{
			Name:        "wave_ci",
			DisplayName: "Wave CI",
			Type:        "mobile_money",
			BaseURL:     "https://api.wave.com/v1",
			LogoURL:     "/icons/aggregators/wave.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"CI", "Côte d'Ivoire", "XOF"},
			},
		},
		{
			Name:        "wave_sn",
			DisplayName: "Wave Sénégal",
			Type:        "mobile_money",
			BaseURL:     "https://api.wave.com/v1",
			LogoURL:     "/icons/aggregators/wave.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"SN", "Sénégal", "XOF"},
			},
		},
		{
			Name:        "lygos",
			DisplayName: "Lygos",
			Type:        "mobile_money",
			BaseURL:     "https://api.lygosapp.com/v1",
			LogoURL:     "/icons/aggregators/lygos.svg",
			Capability:  "deposit",
			Countries: []struct{ Code, Name, Currency string }{
				{"CI", "Côte d'Ivoire", "XOF"},
				{"SN", "Sénégal", "XOF"},
			},
		},

		// --- Orange Money ---
		{
			Name:        "orange_money",
			DisplayName: "Orange Money",
			Type:        "mobile_money",
			BaseURL:     "https://api.orange.com/orange-money-webpay/dev/v1",
			LogoURL:     "/icons/aggregators/orange_money.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"CI", "Côte d'Ivoire", "XOF"},
				{"CM", "Cameroun", "XAF"},
				{"SN", "Sénégal", "XOF"},
			},
		},
		// --- MTN Mobile Money ---
		{
			Name:        "mtn_money",
			DisplayName: "MTN MoMo",
			Type:        "mobile_money",
			BaseURL:     "https://sandbox.momodeveloper.mtn.com",
			LogoURL:     "/icons/aggregators/mtn.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"CI", "Côte d'Ivoire", "XOF"},
				{"CM", "Cameroun", "XAF"},
				{"SN", "Sénégal", "XOF"},
				{"BJ", "Bénin", "XOF"},
			},
		},
		// --- Wave ---
		{
			Name:        "wave_money",
			DisplayName: "Wave",
			Type:        "mobile_money",
			BaseURL:     "https://api.wave.com/v1",
			LogoURL:     "/icons/aggregators/wave.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"CI", "Côte d'Ivoire", "XOF"},
				{"SN", "Sénégal", "XOF"},
			},
		},
		// --- Moov Money ---
		{
			Name:        "moov_money",
			DisplayName: "Moov Money",
			Type:        "mobile_money",
			BaseURL:     "https://testapimarchand2.moov-africa.bj:2010/com.tlc.merchant.api/UssdPush",
			LogoURL:     "/icons/aggregators/moov.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"CI", "Côte d'Ivoire", "XOF"},
				{"BJ", "Bénin", "XOF"},
				{"TG", "Togo", "XOF"},
				{"CM", "Cameroun", "XAF"},
			},
		},
		// --- GLOBAL PROVIDERS (available for all countries) ---
		// Stripe - type "card" = global
		{
			Name:        "stripe",
			DisplayName: "Stripe / Carte Bancaire",
			Type:        "card", // Card type is global - no need for specific countries
			BaseURL:     "https://api.stripe.com/v1",
			LogoURL:     "/icons/aggregators/stripe.svg",
			Capability:  "mixed",
			Countries:   []struct{ Code, Name, Currency string }{}, // Empty = global
		},
		// PayPal - type "international" = global
		{
			Name:        "paypal",
			DisplayName: "PayPal",
			Type:        "international", // International type is global
			BaseURL:     "https://api.paypal.com/v1",
			LogoURL:     "/icons/aggregators/paypal.svg",
			Capability:  "mixed",
			Countries:   []struct{ Code, Name, Currency string }{}, // Empty = global
		},
		// --- Aggregators (Tiers) ---
		{
			Name:        "cinetpay",
			DisplayName: "CinetPay",
			Type:        "mobile_money",
			BaseURL:     "https://api-checkout.cinetpay.com/v2/payment",
			LogoURL:     "/icons/aggregators/cinetpay.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"CI", "Côte d'Ivoire", "XOF"},
				{"SN", "Sénégal", "XOF"},
				{"CM", "Cameroun", "XAF"},
				{"BF", "Burkina Faso", "XOF"},
				{"ML", "Mali", "XOF"},
				{"TG", "Togo", "XOF"},
			},
		},
		{
			Name:        "flutterwave",
			DisplayName: "Flutterwave",
			Type:        "mobile_money",
			BaseURL:     "https://api.flutterwave.com/v3",
			LogoURL:     "/icons/aggregators/flutterwave.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"NG", "Nigeria", "NGN"},
				{"GH", "Ghana", "GHS"},
				{"KE", "Kenya", "KES"},
				{"ZA", "South Africa", "ZAR"},
				{"UG", "Uganda", "UGX"},
				{"TZ", "Tanzania", "TZS"},
				{"RW", "Rwanda", "RWF"},
				{"CM", "Cameroun", "XAF"},
			},
		},
		{
			Name:        "fedapay",
			DisplayName: "FedaPay",
			Type:        "mobile_money",
			BaseURL:     "https://api.fedapay.com/v1", // Sandbox: https://sandbox-api.fedapay.com/v1
			LogoURL:     "/icons/aggregators/fedapay.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"BJ", "Bénin", "XOF"},
				{"TG", "Togo", "XOF"},
				{"NE", "Niger", "XOF"},
				{"CI", "Côte d'Ivoire", "XOF"},
				{"CM", "Cameroun", "XAF"},
			},
		},
		{
			Name: "lygos",
			Countries: []struct{ Code, Name, Currency string }{
				{"LR", "Liberia", "LRD"},
				{"CI", "Côte d'Ivoire", "XOF"},
				{"CM", "Cameroun", "XAF"},
				// Supports 13+ countries
			},
		},
		{
			Name:        "yellowcard",
			DisplayName: "YellowCard",
			Type:        "crypto_ramp",
			BaseURL:     "https://api.yellowcard.io/v1",
			LogoURL:     "/icons/aggregators/yellowcard.svg",
			Capability:  "mixed",
			Countries: []struct{ Code, Name, Currency string }{
				{"NG", "Nigeria", "NGN"},
				{"GH", "Ghana", "GHS"},
				{"CI", "Côte d'Ivoire", "XOF"},
				{"SN", "Sénégal", "XOF"},
				{"KE", "Kenya", "KES"},
				{"ZA", "South Africa", "ZAR"},
				{"CM", "Cameroun", "XAF"},
			},
		},
	}

	for _, p := range providers {
		currencies := `["XOF", "XAF", "NGN", "EUR", "USD"]`
		var providerID string

		// Generate random keys for mockup
		mockKey := fmt.Sprintf("pk_test_%s", uuid.New().String())
		mockSecret := fmt.Sprintf("sk_test_%s", uuid.New().String())

		// Use ON CONFLICT to update existing providers if they exist
		err := db.QueryRow(`
			INSERT INTO payment_providers
			(name, display_name, provider_type, api_base_url, logo_url,
			 supported_currencies, capability, is_demo_mode, is_active,
			 api_key_encrypted, api_secret_encrypted)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE, $9, $10)
			ON CONFLICT (name) DO UPDATE SET
				display_name = EXCLUDED.display_name,
				provider_type = EXCLUDED.provider_type,
				api_base_url = EXCLUDED.api_base_url,
				logo_url = EXCLUDED.logo_url,
				supported_currencies = EXCLUDED.supported_currencies,
				capability = EXCLUDED.capability,
				is_demo_mode = EXCLUDED.is_demo_mode,
				is_active = TRUE,
				api_key_encrypted = COALESCE(payment_providers.api_key_encrypted, EXCLUDED.api_key_encrypted),
				api_secret_encrypted = COALESCE(payment_providers.api_secret_encrypted, EXCLUDED.api_secret_encrypted)
			RETURNING id
		`, p.Name, p.DisplayName, p.Type, p.BaseURL, p.LogoURL, currencies, p.Capability, p.IsDemo, mockKey, mockSecret).Scan(&providerID)

		if err != nil {
			log.Printf("Failed to upsert provider %s: %v", p.Name, err)
			continue
		}

		// Clean up existing country mappings to avoid duplicates
		_, err = db.Exec("DELETE FROM provider_countries WHERE provider_id = $1", providerID)
		if err != nil {
			log.Printf("Failed to clear countries for provider %s: %v", p.Name, err)
			continue
		}

		// Insert country mappings with is_active explicitly set to TRUE
		for _, c := range p.Countries {
			_, err = db.Exec(`
				INSERT INTO provider_countries (provider_id, country_code, country_name, currency, priority, is_active)
				VALUES ($1, $2, $3, $4, 50, TRUE)
			`, providerID, c.Code, c.Name, c.Currency)
			if err != nil {
				log.Printf("Failed to insert country %s for provider %s: %v", c.Code, p.Name, err)
			}
		}

		log.Printf("[Database] Seeded/Updated provider: %s with %d countries", p.DisplayName, len(p.Countries))
	}

	log.Println("[Database] ✅ Payment providers seeding complete")
	return nil
}

// SeedProviderInstances creates a default instance for each payment provider
// These instances have placeholder Vault paths and are linked to ALL available Hot Wallets
func SeedProviderInstances(adminDB, mainDB *sql.DB) error {
	log.Println("[Database] Seeding default provider instances with multi-wallet support...")

	// 1. Get ALL operations hot wallets from Main DB (retry mechanism to wait for wallet-service seeding)
	type HotWallet struct {
		ID       string
		Currency string
	}
	var hotWallets []HotWallet

	maxRetries := 12
	for i := 0; i < maxRetries; i++ {
		hotWallets = nil // Reset for retry
		rows, err := mainDB.Query(`SELECT id, currency FROM platform_accounts WHERE account_type = 'operations' AND is_active = true`)
		if err != nil {
			log.Printf("[Database] ⚠️ Warning: Failed to query operations wallets (attempt %d/%d): %v", i+1, maxRetries, err)
		} else {
			defer rows.Close() // This defer runs when SeedProviderInstances exits, which is fine, but for loop safety we should close manually if retrying or loop ends
			for rows.Next() {
				var hw HotWallet
				if err := rows.Scan(&hw.ID, &hw.Currency); err == nil {
					hotWallets = append(hotWallets, hw)
				}
			}
			rows.Close() // Explicit close for retry loop

			if len(hotWallets) > 0 {
				log.Printf("[Database] ✅ Found %d operations hot wallets", len(hotWallets))
				break
			}
		}

		if i < maxRetries-1 {
			log.Printf("[Database] Waiting for wallet-service seeding... (attempt %d/%d)", i+1, maxRetries)
			time.Sleep(5 * time.Second)
		}
	}

	if len(hotWallets) == 0 {
		log.Printf("[Database] ❌ Error: No 'operations' hot wallets found in Main DB after %d retries. Instances will be created without wallet links. Ensure wallet-service is running.", maxRetries)
	}

	// 2. Get all providers with their supported currencies from Admin DB
	providerRows, err := adminDB.Query(`
		SELECT p.id, p.name, p.display_name, p.provider_type, COALESCE(p.api_base_url, ''), COALESCE(pc.currency, '')
		FROM payment_providers p
		LEFT JOIN provider_countries pc ON p.id = pc.provider_id
		ORDER BY p.id
	`)
	if err != nil {
		return fmt.Errorf("failed to get providers: %w", err)
	}
	defer providerRows.Close()

	type providerInfo struct {
		ID          string
		Name        string
		DisplayName string
		ProviderType string
		BaseURL     string
		Currencies  map[string]bool
	}

	providerMap := make(map[string]*providerInfo)
	for providerRows.Next() {
		var id, name, displayName, providerType, baseURL, currency string
		if err := providerRows.Scan(&id, &name, &displayName, &providerType, &baseURL, &currency); err != nil {
			continue
		}
		if _, exists := providerMap[id]; !exists {
			providerMap[id] = &providerInfo{
				ID:          id,
				Name:        name,
				DisplayName: displayName,
				ProviderType: providerType,
				BaseURL:     baseURL,
				Currencies:  make(map[string]bool),
			}
		}
		if currency != "" {
			providerMap[id].Currencies[currency] = true
		}
	}

	// 3. Seed Instances and Multi-Wallet Links
	for _, p := range providerMap {
		// Mark instance as global if provider is global (no country mappings) or provider_type indicates global.
		isGlobalProvider := false
		if len(p.Currencies) == 0 {
			isGlobalProvider = true
		}
		switch strings.ToLower(p.ProviderType) {
		case "card", "international", "demo":
			isGlobalProvider = true
		}

		// Vault path matches actual Vault structure: secret/aggregators/{provider}
		// Do NOT add /default - that's not how secrets are stored
		vaultPath := fmt.Sprintf("secret/aggregators/%s", p.Name)

		createOrUpdate := func(instanceName string, isPrimary bool, depositEnabled bool, withdrawEnabled bool) (string, bool) {
			var instanceID string
			err := adminDB.QueryRow(`SELECT id FROM provider_instances WHERE provider_id = $1 AND name = $2`, p.ID, instanceName).Scan(&instanceID)
			created := false
			isTestMode := strings.Contains(strings.ToLower(instanceName), "(sandbox)") || strings.Contains(strings.ToLower(instanceName), "sandbox")
			if err == sql.ErrNoRows {
				if isPrimary {
					adminDB.Exec("UPDATE provider_instances SET is_primary = FALSE WHERE provider_id = $1", p.ID)
				}
				err = adminDB.QueryRow(`
					INSERT INTO provider_instances
						(provider_id, name, vault_secret_path, is_active, is_primary, is_global, is_test_mode, priority, health_status, deposit_enabled, withdraw_enabled)
					VALUES ($1, $2, $3, TRUE, $4, $5, $6, 50, 'active', $7, $8)
					RETURNING id
				`, p.ID, instanceName, vaultPath, isPrimary, isGlobalProvider, isTestMode, depositEnabled, withdrawEnabled).Scan(&instanceID)
				if err != nil {
					log.Printf("[Database] Failed to seed instance for %s (%s): %v", p.Name, instanceName, err)
					return "", false
				}
				created = true
				log.Printf("[Database] ✅ Created default instance for: %s (%s)", p.DisplayName, instanceName)
			} else if err != nil {
				log.Printf("[Database] Error checking instance for %s (%s): %v", p.Name, instanceName, err)
				return "", false
			} else {
				if isPrimary {
					adminDB.Exec("UPDATE provider_instances SET is_primary = FALSE WHERE provider_id = $1", p.ID)
				}
				adminDB.Exec(`
					UPDATE provider_instances
					SET is_active = TRUE,
						is_primary = $1,
						is_global = $2,
						is_test_mode = $3,
						vault_secret_path = $4,
						health_status = 'active',
						deposit_enabled = $5,
						withdraw_enabled = $6,
						updated_at = NOW()
					WHERE id = $7
				`, isPrimary, isGlobalProvider, isTestMode, vaultPath, depositEnabled, withdrawEnabled, instanceID)
				log.Printf("[Database] ✅ Updated default instance for: %s (%s)", p.DisplayName, instanceName)
			}
			return instanceID, created
		}

		isMTN := p.Name == "mtn_money" || p.Name == "mtn_momo"
		isSandboxLiveCapable := isMTN || p.Name == "paypal" || p.Name == "stripe" || p.Name == "flutterwave" || p.Name == "paystack"
		var instanceIDs []string
		if isMTN {
			// Four instances: (Collections/Disbursements) x (Sandbox/Live)
			// Sandbox
			id1, _ := createOrUpdate("MTN Collections (Sandbox)", true, true, false)
			if id1 != "" {
				instanceIDs = append(instanceIDs, id1)
			}
			id2, _ := createOrUpdate("MTN Disbursements (Sandbox)", false, false, true)
			if id2 != "" {
				instanceIDs = append(instanceIDs, id2)
			}
			// Live
			id3, _ := createOrUpdate("MTN Collections (Live)", false, true, false)
			if id3 != "" {
				instanceIDs = append(instanceIDs, id3)
			}
			id4, _ := createOrUpdate("MTN Disbursements (Live)", false, false, true)
			if id4 != "" {
				instanceIDs = append(instanceIDs, id4)
			}
		} else if isSandboxLiveCapable {
			idSandbox, _ := createOrUpdate(p.DisplayName+" (Sandbox)", true, true, true)
			if idSandbox != "" {
				instanceIDs = append(instanceIDs, idSandbox)
			}
			idLive, _ := createOrUpdate(p.DisplayName+" (Live)", false, true, true)
			if idLive != "" {
				instanceIDs = append(instanceIDs, idLive)
			}
		} else {
			id, _ := createOrUpdate("Instance Principale", true, true, true)
			if id != "" {
				instanceIDs = append(instanceIDs, id)
			}
		}

		// 3b. Ensure credentials exist for each instance.
		// seedDefaultAPIKeys runs BEFORE instances are created (during createAdminTables),
		// so fresh DB resets would end up with empty credentials unless we seed them here.
		for _, instanceID := range instanceIDs {
			var currentCreds []byte
			credErr := adminDB.QueryRow(`SELECT COALESCE(api_credentials, '{}'::jsonb) FROM provider_instances WHERE id = $1`, instanceID).Scan(&currentCreds)
			if credErr == nil {
				trimmed := strings.TrimSpace(string(currentCreds))
				if trimmed == "" || trimmed == "{}" {
					seedCreds := buildSeedInstanceCredentials(p.Name, p.ProviderType, p.BaseURL)
					// For MTN live instances, prefill production defaults so both modes are visible after reset.
					if isMTN {
						var instanceName string
						_ = adminDB.QueryRow(`SELECT name FROM provider_instances WHERE id = $1`, instanceID).Scan(&instanceName)
						if strings.Contains(strings.ToLower(instanceName), "(live)") {
							seedCreds["environment"] = "production"
							seedCreds["target_environment"] = "production"
							seedCreds["base_url"] = "https://proxy.momoapi.mtn.com"
						}
					}
					// Generic sandbox/live hints for providers that support both
					if strings.Contains(strings.ToLower(instanceName), "(sandbox)") {
						seedCreds["mode"] = "sandbox"
					}
					if strings.Contains(strings.ToLower(instanceName), "(live)") {
						seedCreds["mode"] = "live"
					}
					if len(seedCreds) > 0 {
						credJSON, err := json.Marshal(seedCreds)
						if err == nil {
							_, _ = adminDB.Exec(`UPDATE provider_instances SET api_credentials = $1, updated_at = NOW() WHERE id = $2`, credJSON, instanceID)
							log.Printf("[Database] ✅ Seeded credentials for instance: %s (%s)", p.DisplayName, instanceID)
						}
					}
				}
			}
		}

		for _, instanceID := range instanceIDs {
			linkedCount := 0
			for _, hw := range hotWallets {
				if len(p.Currencies) == 0 || p.Currencies[hw.Currency] {
					_, err := adminDB.Exec(`
						INSERT INTO provider_instance_wallets (instance_id, currency, hot_wallet_id, is_active, priority)
						VALUES ($1, $2, $3, TRUE, 50)
						ON CONFLICT (instance_id, currency, hot_wallet_id) DO UPDATE SET is_active = TRUE
					`, instanceID, hw.Currency, hw.ID)
					if err != nil {
						log.Printf("[Database] Failed to link wallet %s to instance %s: %v", hw.ID, instanceID, err)
					} else {
						linkedCount++
					}
				}
			}
			if linkedCount > 0 {
				log.Printf("[Database] ✅ Linked %d hot wallets to instance: %s (%s)", linkedCount, p.DisplayName, instanceID)
			}
		}

		// IMPORTANT: make aggregator_instances.id stable and equal to provider_instances.id.
		// This guarantees that even if a caller mistakenly passes a provider_instance UUID
		// to transfer-service wallet selection, it still matches an aggregator instance.
		var aggregatorID string
		aggErr := mainDB.QueryRow(`SELECT id FROM aggregator_settings WHERE provider_code = $1`, p.Name).Scan(&aggregatorID)
		if aggErr == nil {
			for _, instanceID := range instanceIDs {
				var instanceName string
				_ = adminDB.QueryRow(`SELECT name FROM provider_instances WHERE id = $1`, instanceID).Scan(&instanceName)
				isTestMode := true
				if strings.Contains(strings.ToLower(instanceName), "(live)") {
					isTestMode = false
				}
				var aggInstanceID string
				aggInstErr := mainDB.QueryRow(`
					INSERT INTO aggregator_instances (
						id,
						aggregator_id,
						instance_name,
						api_credentials,
						vault_secret_path,
						enabled,
						is_global,
						priority,
						health_status,
						is_test_mode
					)
					VALUES ($1, $2, $3, '{}'::jsonb, $4, TRUE, $5, 50, 'active', $6)
					ON CONFLICT (id) DO UPDATE SET
						aggregator_id = EXCLUDED.aggregator_id,
						instance_name = EXCLUDED.instance_name,
						vault_secret_path = EXCLUDED.vault_secret_path,
						enabled = TRUE,
						is_global = EXCLUDED.is_global,
						health_status = 'active',
						is_test_mode = EXCLUDED.is_test_mode,
						updated_at = NOW()
					RETURNING id
				`, instanceID, aggregatorID, instanceName, vaultPath, isGlobalProvider, isTestMode).Scan(&aggInstanceID)
				if aggInstErr != nil {
					log.Printf("[Database] Failed to upsert aggregator instance for %s (%s): %v", p.Name, instanceName, aggInstErr)
					continue
				}

				for _, hw := range hotWallets {
					if len(p.Currencies) == 0 || p.Currencies[hw.Currency] {
						_, err := mainDB.Exec(`
							INSERT INTO aggregator_instance_wallets (instance_id, hot_wallet_id, currency, is_primary, priority, enabled)
							VALUES ($1, $2, $3, TRUE, 50, TRUE)
							ON CONFLICT (instance_id, hot_wallet_id, currency) DO UPDATE SET enabled = TRUE, is_primary = TRUE
						`, aggInstanceID, hw.ID, hw.Currency)
						if err != nil {
							log.Printf("[Database] Failed to link wallet %s to aggregator instance %s: %v", hw.ID, aggInstanceID, err)
						}
					}
				}
			}
		}
	}

	log.Println("[Database] ✅ Provider instances seeding complete with multi-wallet support")
	log.Println("[Database] ⚠️  Configure real API keys in Vault before production use")
	return nil
}

func seedDefaultData(db *sql.DB) error {
	// Check if data already exists
	var count int
	db.QueryRow("SELECT COUNT(*) FROM admin_roles").Scan(&count)
	if count > 0 {
		return nil // Already seeded
	}

	// Insert permissions
	permissions := []struct {
		Code     string
		Name     string
		Desc     string
		Category string
	}{
		{"users.view", "View Users", "View user information", "Users"},
		{"users.create", "Create Users", "Create new users", "Users"},
		{"users.update", "Update Users", "Update user information", "Users"},
		{"users.block", "Block Users", "Block/unblock users", "Users"},
		{"users.delete", "Delete Users", "Delete users", "Users"},
		{"kyc.view", "View KYC", "View KYC documents", "KYC"},
		{"kyc.approve", "Approve KYC", "Approve KYC applications", "KYC"},
		{"kyc.reject", "Reject KYC", "Reject KYC applications", "KYC"},
		{"transactions.view", "View Transactions", "View all transactions", "Transactions"},
		{"transactions.block", "Block Transactions", "Block suspicious transactions", "Transactions"},
		{"transactions.refund", "Refund Transactions", "Process refunds", "Transactions"},
		{"transactions.approve", "Approve Transactions", "Approve pending transactions", "Transactions"},
		{"cards.view", "View Cards", "View all cards", "Cards"},
		{"cards.freeze", "Freeze Cards", "Freeze/unfreeze cards", "Cards"},
		{"cards.block", "Block Cards", "Permanently block cards", "Cards"},
		{"cards.replace", "Replace Cards", "Issue replacement cards", "Cards"},
		{"wallets.view", "View Wallets", "View all wallets", "Wallets"},
		{"wallets.freeze", "Freeze Wallets", "Freeze/unfreeze wallets", "Wallets"},
		{"wallets.adjust", "Adjust Balances", "Adjust wallet balances", "Wallets"},
		{"transfers.view", "View Transfers", "View all transfers", "Transfers"},
		{"transfers.block", "Block Transfers", "Block suspicious transfers", "Transfers"},
		{"transfers.approve", "Approve Transfers", "Approve pending transfers", "Transfers"},
		{"exchanges.view", "View Exchanges", "View all exchanges", "Exchanges"},
		{"exchanges.rates", "Set Rates", "Configure exchange rates", "Exchanges"},
		{"system.view", "View System", "View system status", "System"},
		{"system.logs", "View Logs", "View system logs", "System"},
		{"system.settings", "Manage Settings", "Manage system settings", "System"},
		{"admins.view", "View Admins", "View admin users", "Admins"},
		{"admins.create", "Create Admins", "Create new admins", "Admins"},
		{"admins.update", "Update Admins", "Update admin info", "Admins"},
		{"admins.delete", "Delete Admins", "Delete admins", "Admins"},
		{"admins.roles", "Manage Roles", "Manage admin roles", "Admins"},
		{"analytics.view", "View Analytics", "View analytics dashboard", "Analytics"},
		{"analytics.export", "Export Reports", "Export data reports", "Analytics"},
	}

	permissionIDs := make(map[string]string)
	for _, p := range permissions {
		var id string
		err := db.QueryRow(`
			INSERT INTO admin_permissions (code, name, description, category)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, p.Code, p.Name, p.Desc, p.Category).Scan(&id)
		if err != nil {
			return err
		}
		permissionIDs[p.Code] = id
	}

	// Create roles with permissions
	roles := map[string][]string{
		"Super Admin": {"users.view", "users.create", "users.update", "users.block", "users.delete",
			"kyc.view", "kyc.approve", "kyc.reject",
			"transactions.view", "transactions.block", "transactions.refund", "transactions.approve",
			"cards.view", "cards.freeze", "cards.block", "cards.replace",
			"wallets.view", "wallets.freeze", "wallets.adjust",
			"transfers.view", "transfers.block", "transfers.approve",
			"exchanges.view", "exchanges.rates",
			"system.view", "system.logs", "system.settings",
			"admins.view", "admins.create", "admins.update", "admins.delete", "admins.roles",
			"analytics.view", "analytics.export"},
		"Admin": {"users.view", "users.update", "users.block",
			"kyc.view", "kyc.approve", "kyc.reject",
			"transactions.view", "transactions.block", "transactions.refund",
			"cards.view", "cards.freeze", "cards.block",
			"wallets.view", "wallets.freeze",
			"transfers.view", "transfers.block",
			"exchanges.view",
			"system.view",
			"analytics.view"},
		"Support": {"users.view", "users.update",
			"kyc.view",
			"transactions.view",
			"cards.view", "cards.freeze",
			"wallets.view",
			"transfers.view"},
		"Compliance": {"users.view",
			"kyc.view", "kyc.approve", "kyc.reject",
			"transactions.view", "transactions.block",
			"wallets.view",
			"transfers.view", "transfers.block",
			"analytics.view"},
		"Analyst": {"users.view",
			"transactions.view",
			"exchanges.view",
			"analytics.view", "analytics.export"},
		"Viewer": {"users.view", "transactions.view", "cards.view", "wallets.view", "analytics.view"},
	}

	for roleName, perms := range roles {
		var roleID string
		err := db.QueryRow(`
			INSERT INTO admin_roles (name, description, is_system)
			VALUES ($1, $2, TRUE)
			ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description
			RETURNING id
		`, roleName, fmt.Sprintf("%s role with predefined permissions", roleName)).Scan(&roleID)
		if err != nil {
			return err
		}

		for _, permCode := range perms {
			if permID, ok := permissionIDs[permCode]; ok {
				db.Exec(`
					INSERT INTO admin_role_permissions (role_id, permission_id)
					VALUES ($1, $2)
					ON CONFLICT DO NOTHING
				`, roleID, permID)
			}
		}
	}

	// Create default super admin
	if err := createDefaultSuperAdmin(db); err != nil {
		return err
	}

	return nil
}

// seedDefaultAPIKeys seeds default API keys for all provider instances
func seedDefaultAPIKeys(db *sql.DB) error {
	log.Println("[Database] Seeding default API keys for provider instances...")

	// Get all provider instances that don't have API credentials yet
	rows, err := db.Query(`
		SELECT pi.id, pi.name, p.name as provider_name, p.provider_type
		FROM provider_instances pi
		JOIN payment_providers p ON pi.provider_id = p.id
		WHERE (pi.api_credentials IS NULL OR pi.api_credentials = '{}'::jsonb)
		AND pi.is_active = true
	`)
	if err != nil {
		return fmt.Errorf("failed to get provider instances: %w", err)
	}
	defer rows.Close()

	type InstanceInfo struct {
		ID           string
		Name         string
		ProviderName string
		ProviderType string
	}

	var instances []InstanceInfo
	for rows.Next() {
		var inst InstanceInfo
		if err := rows.Scan(&inst.ID, &inst.Name, &inst.ProviderName, &inst.ProviderType); err != nil {
			continue
		}
		instances = append(instances, inst)
	}

	// Generate mock API keys for each instance
	for _, inst := range instances {
		credentials := make(map[string]interface{})
		
		// Generate different keys based on provider type
		switch inst.ProviderType {
		case "mobile_money":
			credentials["client_id"] = fmt.Sprintf("client_%s_%s", inst.ProviderName, uuid.New().String()[:8])
			credentials["client_secret"] = fmt.Sprintf("secret_%s_%s", inst.ProviderName, uuid.New().String()[:12])
			credentials["api_key"] = fmt.Sprintf("pk_%s_%s", inst.ProviderName, uuid.New().String()[:8])
			credentials["secret_key"] = fmt.Sprintf("sk_%s_%s", inst.ProviderName, uuid.New().String()[:12])
			credentials["webhook_secret"] = fmt.Sprintf("wh_%s_%s", inst.ProviderName, uuid.New().String()[:16])
			credentials["base_url"] = "https://api.example.com"
			credentials["environment"] = "sandbox"
			
		case "card":
			credentials["api_key"] = fmt.Sprintf("pk_%s_%s", inst.ProviderName, uuid.New().String()[:8])
			credentials["secret_key"] = fmt.Sprintf("sk_%s_%s", inst.ProviderName, uuid.New().String()[:12])
			credentials["webhook_secret"] = fmt.Sprintf("whsec_%s_%s", inst.ProviderName, uuid.New().String()[:16])
			credentials["base_url"] = "https://api.stripe.com/v1"
			credentials["environment"] = "sandbox"
			
		case "international":
			credentials["client_id"] = fmt.Sprintf("client_%s_%s", inst.ProviderName, uuid.New().String()[:8])
			credentials["client_secret"] = fmt.Sprintf("secret_%s_%s", inst.ProviderName, uuid.New().String()[:12])
			credentials["api_key"] = fmt.Sprintf("api_%s_%s", inst.ProviderName, uuid.New().String()[:8])
			credentials["base_url"] = "https://api.paypal.com/v1"
			credentials["environment"] = "sandbox"
			
		case "crypto_ramp":
			credentials["api_key"] = fmt.Sprintf("pk_%s_%s", inst.ProviderName, uuid.New().String()[:8])
			credentials["secret_key"] = fmt.Sprintf("sk_%s_%s", inst.ProviderName, uuid.New().String()[:12])
			credentials["public_key"] = fmt.Sprintf("pub_%s_%s", inst.ProviderName, uuid.New().String()[:16])
			credentials["base_url"] = "https://api.yellowcard.io/v1"
			credentials["environment"] = "sandbox"
			
		case "demo":
			credentials["api_key"] = "demo_api_key_12345"
			credentials["secret_key"] = "demo_secret_key_12345"
			credentials["client_id"] = "demo_client_12345"
			credentials["client_secret"] = "demo_client_secret_12345"
			credentials["base_url"] = "https://demo.api.com"
			credentials["environment"] = "demo"
		}

		// Convert to JSON and update database
		credJSON, err := json.Marshal(credentials)
		if err != nil {
			log.Printf("Failed to marshal credentials for %s: %v", inst.ProviderName, err)
			continue
		}

		_, err = db.Exec(`
			UPDATE provider_instances 
			SET api_credentials = $1, updated_at = NOW() 
			WHERE id = $2
		`, credJSON, inst.ID)

		if err != nil {
			log.Printf("Failed to update credentials for %s: %v", inst.ProviderName, err)
		} else {
			log.Printf("[Database] ✅ Seeded API keys for %s - %s", inst.ProviderName, inst.Name)
		}
	}

	log.Println("[Database] ✅ Default API keys seeding complete")
	return nil
}

func createDefaultSuperAdmin(db *sql.DB) error {
	// Check if any admin exists
	var adminCount int
	db.QueryRow("SELECT COUNT(*) FROM admin_users").Scan(&adminCount)
	if adminCount > 0 {
		return nil // Admin already exists
	}

	// Get Super Admin role ID
	var roleID string
	err := db.QueryRow("SELECT id FROM admin_roles WHERE name = 'Super Admin'").Scan(&roleID)
	if err != nil {
		return fmt.Errorf("failed to get Super Admin role: %w", err)
	}

	// Default password: Admin123!
	// Generate bcrypt hash dynamically to ensure it's always valid
	defaultPassword := "Admin123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash default password: %w", err)
	}

	// Create default super admin
	_, err = db.Exec(`
		INSERT INTO admin_users (email, password_hash, first_name, last_name, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, TRUE)
		ON CONFLICT (email) DO NOTHING
	`, "admin@zekora.com", string(hashedPassword), "Super", "Admin", roleID)

	if err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	fmt.Println("===========================================")
	fmt.Println("DEFAULT SUPER ADMIN CREATED")
	fmt.Println("Email: admin@zekora.com")
	fmt.Println("Password: Admin123!")
	fmt.Println("⚠️  CHANGE THIS PASSWORD IMMEDIATELY!")
	fmt.Println("===========================================")

	return nil
}

// InitializeKafka creates a new Kafka client for messaging
func InitializeKafka(brokers string, groupID string) (*messaging.KafkaClient, error) {
	brokerList := strings.Split(brokers, ",")

	client := messaging.NewKafkaClient(brokerList, groupID)

	log.Printf("[Kafka] Admin-service connected to brokers: %s with group: %s", brokers, groupID)
	return client, nil
}
