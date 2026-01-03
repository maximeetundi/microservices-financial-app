package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "crypto_bank")
	sslmode := getEnv("DB_SSL_MODE", "disable")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	return nil
}

func InitSchema() error {
	schema := `
	-- Events table
	CREATE TABLE IF NOT EXISTS events (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		creator_id UUID NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		location TEXT,
		cover_image TEXT,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP NOT NULL,
		sale_start_date TIMESTAMP NOT NULL,
		sale_end_date TIMESTAMP NOT NULL,
		form_fields JSONB DEFAULT '[]',
		qr_code TEXT,
		event_code VARCHAR(20) UNIQUE NOT NULL,
		status VARCHAR(20) DEFAULT 'draft',
		currency VARCHAR(10) DEFAULT 'XOF',
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	-- Ticket tiers table
	CREATE TABLE IF NOT EXISTS ticket_tiers (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		event_id UUID REFERENCES events(id) ON DELETE CASCADE,
		name VARCHAR(100) NOT NULL,
		icon VARCHAR(50) DEFAULT '🎫',
		price DECIMAL(15,2) NOT NULL,
		quantity INT DEFAULT -1,
		sold INT DEFAULT 0,
		description TEXT,
		benefits JSONB DEFAULT '[]',
		color VARCHAR(20) DEFAULT '#6366f1',
		sort_order INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT NOW()
	);

	-- Tickets table (purchased)
	CREATE TABLE IF NOT EXISTS tickets (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		event_id UUID REFERENCES events(id),
		buyer_id UUID NOT NULL,
		tier_id UUID REFERENCES ticket_tiers(id),
		tier_name VARCHAR(100) NOT NULL,
		tier_icon VARCHAR(50) DEFAULT '🎫',
		price DECIMAL(15,2) NOT NULL,
		currency VARCHAR(10) NOT NULL,
		form_data JSONB DEFAULT '{}',
		qr_code TEXT,
		ticket_code VARCHAR(30) UNIQUE NOT NULL,
		status VARCHAR(20) DEFAULT 'pending',
		transaction_id UUID,
		used_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT NOW()
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_events_creator ON events(creator_id);
	CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);
	CREATE INDEX IF NOT EXISTS idx_events_code ON events(event_code);
	CREATE INDEX IF NOT EXISTS idx_ticket_tiers_event ON ticket_tiers(event_id);
	CREATE INDEX IF NOT EXISTS idx_tickets_event ON tickets(event_id);
	CREATE INDEX IF NOT EXISTS idx_tickets_buyer ON tickets(buyer_id);
	CREATE INDEX IF NOT EXISTS idx_tickets_code ON tickets(ticket_code);
	CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Database schema initialized")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
