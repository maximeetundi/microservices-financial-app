package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

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
		
		// Create index for faster queries
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_admin_id ON admin_audit_logs(admin_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON admin_audit_logs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_admin_sessions_admin_id ON admin_sessions(admin_id)`,
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

	return nil
}

// RabbitMQ Client
type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func InitializeRabbitMQ(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare admin exchange
	err = ch.ExchangeDeclare("admin.events", "topic", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare admin command exchange for sending commands to other services
	err = ch.ExchangeDeclare("admin.commands", "topic", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare command exchange: %w", err)
	}

	return &RabbitMQClient{conn: conn, channel: ch}, nil
}

func (r *RabbitMQClient) GetChannel() *amqp.Channel {
	return r.channel
}

func (r *RabbitMQClient) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *RabbitMQClient) Publish(exchange, routingKey string, message []byte) error {
	return r.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}
