package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// InitDB initializes the database connection and creates tables
func InitDB() (*sql.DB, error) {
	// Read from environment variables or use defaults
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "postgres"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "admin"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "secure_password"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "crypto_bank"
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database successfully")

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS associations (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		description TEXT,
		rules JSONB,
		total_members INT DEFAULT 0,
		treasury_balance DECIMAL(15,2) DEFAULT 0,
		currency VARCHAR(3) DEFAULT 'XOF',
		status VARCHAR(50) DEFAULT 'active',
		created_by VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS members (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		association_id UUID NOT NULL REFERENCES associations(id) ON DELETE CASCADE,
		user_id VARCHAR(255) NOT NULL,
		role VARCHAR(50) DEFAULT 'member',
		join_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status VARCHAR(50) DEFAULT 'pending',
		contributions_paid DECIMAL(15,2) DEFAULT 0,
		loans_received DECIMAL(15,2) DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(association_id, user_id)
	);

	CREATE TABLE IF NOT EXISTS meetings (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		association_id UUID NOT NULL REFERENCES associations(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		date TIMESTAMP NOT NULL,
		location VARCHAR(255),
		agenda JSONB,
		minutes TEXT,
		attendance JSONB,
		status VARCHAR(50) DEFAULT 'scheduled',
		created_by VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS treasury_transactions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		association_id UUID NOT NULL REFERENCES associations(id) ON DELETE CASCADE,
		type VARCHAR(50) NOT NULL,
		amount DECIMAL(15,2) NOT NULL,
		from_member_id UUID REFERENCES members(id),
		to_member_id UUID REFERENCES members(id),
		description TEXT,
		status VARCHAR(50) DEFAULT 'pending',
		created_by VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS contributions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		association_id UUID NOT NULL REFERENCES associations(id) ON DELETE CASCADE,
		member_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
		amount DECIMAL(15,2) NOT NULL,
		period VARCHAR(50) NOT NULL,
		due_date TIMESTAMP NOT NULL,
		paid_date TIMESTAMP,
		status VARCHAR(50) DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS loans (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		association_id UUID NOT NULL REFERENCES associations(id) ON DELETE CASCADE,
		borrower_id UUID NOT NULL REFERENCES members(id) ON DELETE CASCADE,
		amount DECIMAL(15,2) NOT NULL,
		interest_rate DECIMAL(5,2) DEFAULT 0,
		duration INT NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		repayments JSONB,
		status VARCHAR(50) DEFAULT 'pending',
		approved_by VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_members_association ON members(association_id);
	CREATE INDEX IF NOT EXISTS idx_members_user ON members(user_id);
	CREATE INDEX IF NOT EXISTS idx_meetings_association ON meetings(association_id);
	CREATE INDEX IF NOT EXISTS idx_treasury_association ON treasury_transactions(association_id);
	CREATE INDEX IF NOT EXISTS idx_contributions_member ON contributions(member_id);
	CREATE INDEX IF NOT EXISTS idx_loans_borrower ON loans(borrower_id);
	`

	// Execute main schema
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	// Add missing columns for existing tables (migrations)
	migrations := []string{
		`ALTER TABLE associations ADD COLUMN IF NOT EXISTS created_by VARCHAR(255)`,
		`UPDATE associations SET created_by = 'system' WHERE created_by IS NULL`,
		`ALTER TABLE associations ALTER COLUMN created_by SET NOT NULL`,
		`ALTER TABLE meetings ADD COLUMN IF NOT EXISTS created_by VARCHAR(255)`,
		`UPDATE meetings SET created_by = 'system' WHERE created_by IS NULL`,
		`ALTER TABLE meetings ALTER COLUMN created_by SET NOT NULL`,
		`ALTER TABLE treasury_transactions ADD COLUMN IF NOT EXISTS created_by VARCHAR(255)`,
		`UPDATE treasury_transactions SET created_by = 'system' WHERE created_by IS NULL`,
		`ALTER TABLE treasury_transactions ALTER COLUMN created_by SET NOT NULL`,
	}

	for _, migration := range migrations {
		_, _ = db.Exec(migration) // Ignore errors as column might already exist
	}

	return nil
}
