package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Initialize(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Create notifications table if not exists
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
		id VARCHAR(36) PRIMARY KEY,
		user_id VARCHAR(36) NOT NULL,
		type VARCHAR(50) NOT NULL,
		title VARCHAR(255) NOT NULL,
		message TEXT NOT NULL,
		data TEXT,
		is_read BOOLEAN DEFAULT FALSE,
		read_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
	CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(created_at DESC);
	`
	_, err := db.Exec(query)
	return err
}
