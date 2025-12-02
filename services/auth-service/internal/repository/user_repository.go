package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.RegisterRequest) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (email, phone, password_hash, first_name, last_name, date_of_birth, country)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, email, phone, first_name, last_name, date_of_birth, country, kyc_status, kyc_level, 
				  is_active, two_fa_enabled, email_verified, phone_verified, created_at, updated_at
	`

	var newUser models.User
	err = r.db.QueryRow(query, user.Email, user.Phone, string(hashedPassword), 
		user.FirstName, user.LastName, user.DateOfBirth, user.Country).Scan(
		&newUser.ID, &newUser.Email, &newUser.Phone, &newUser.FirstName, &newUser.LastName,
		&newUser.DateOfBirth, &newUser.Country, &newUser.KYCStatus, &newUser.KYCLevel,
		&newUser.IsActive, &newUser.TwoFAEnabled, &newUser.EmailVerified, &newUser.PhoneVerified,
		&newUser.CreatedAt, &newUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, phone, password_hash, first_name, last_name, date_of_birth, country,
			   kyc_status, kyc_level, is_active, two_fa_enabled, two_fa_secret, email_verified,
			   phone_verified, last_login_at, created_at, updated_at, failed_attempts, locked_until
		FROM users WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Phone, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.DateOfBirth, &user.Country, &user.KYCStatus, &user.KYCLevel, &user.IsActive,
		&user.TwoFAEnabled, &user.TwoFASecret, &user.EmailVerified, &user.PhoneVerified,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt, &user.FailedAttempts, &user.LockedUntil,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(userID string) (*models.User, error) {
	query := `
		SELECT id, email, phone, first_name, last_name, date_of_birth, country,
			   kyc_status, kyc_level, is_active, two_fa_enabled, email_verified,
			   phone_verified, last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.Phone, &user.FirstName, &user.LastName,
		&user.DateOfBirth, &user.Country, &user.KYCStatus, &user.KYCLevel, &user.IsActive,
		&user.TwoFAEnabled, &user.EmailVerified, &user.PhoneVerified,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (r *UserRepository) UpdateLastLogin(userID string) error {
	query := `UPDATE users SET last_login_at = NOW(), failed_attempts = 0, locked_until = NULL WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *UserRepository) IncrementFailedAttempts(userID string, maxAttempts int, lockoutDuration time.Duration) error {
	query := `
		UPDATE users 
		SET failed_attempts = failed_attempts + 1,
			locked_until = CASE 
				WHEN failed_attempts + 1 >= $2 THEN NOW() + INTERVAL '%d minutes'
				ELSE locked_until 
			END
		WHERE id = $1
	`
	
	formattedQuery := fmt.Sprintf(query, int(lockoutDuration.Minutes()))
	_, err := r.db.Exec(formattedQuery, userID, maxAttempts)
	return err
}

func (r *UserRepository) IsLocked(userID string) (bool, error) {
	query := `SELECT locked_until FROM users WHERE id = $1`
	
	var lockedUntil *time.Time
	err := r.db.QueryRow(query, userID).Scan(&lockedUntil)
	if err != nil {
		return false, err
	}

	if lockedUntil == nil {
		return false, nil
	}

	return time.Now().Before(*lockedUntil), nil
}

func (r *UserRepository) UpdateEmailVerification(userID string, verified bool) error {
	query := `UPDATE users SET email_verified = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, verified, userID)
	return err
}

func (r *UserRepository) UpdatePhoneVerification(userID string, verified bool) error {
	query := `UPDATE users SET phone_verified = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, verified, userID)
	return err
}

func (r *UserRepository) Update2FASecret(userID, secret string) error {
	query := `UPDATE users SET two_fa_secret = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, secret, userID)
	return err
}

func (r *UserRepository) Enable2FA(userID string, enabled bool) error {
	query := `UPDATE users SET two_fa_enabled = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, enabled, userID)
	return err
}

func (r *UserRepository) UpdatePassword(userID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err = r.db.Exec(query, string(hashedPassword), userID)
	return err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) PhoneExists(phone string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`
	var exists bool
	err := r.db.QueryRow(query, phone).Scan(&exists)
	return exists, err
}

// Verification token methods
func (r *UserRepository) CreateVerificationToken(userID, tokenType string, expiresAt time.Time) (*models.VerificationToken, error) {
	query := `
		INSERT INTO verification_tokens (user_id, token, type, expires_at)
		VALUES ($1, gen_random_uuid()::text, $2, $3)
		RETURNING id, user_id, token, type, expires_at, used, created_at
	`

	var token models.VerificationToken
	err := r.db.QueryRow(query, userID, tokenType, expiresAt).Scan(
		&token.ID, &token.UserID, &token.Token, &token.Type,
		&token.ExpiresAt, &token.Used, &token.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create verification token: %w", err)
	}

	return &token, nil
}

func (r *UserRepository) GetVerificationToken(token, tokenType string) (*models.VerificationToken, error) {
	query := `
		SELECT id, user_id, token, type, expires_at, used, created_at
		FROM verification_tokens
		WHERE token = $1 AND type = $2 AND used = false AND expires_at > NOW()
	`

	var verificationToken models.VerificationToken
	err := r.db.QueryRow(query, token, tokenType).Scan(
		&verificationToken.ID, &verificationToken.UserID, &verificationToken.Token,
		&verificationToken.Type, &verificationToken.ExpiresAt, &verificationToken.Used,
		&verificationToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("verification token not found or expired")
		}
		return nil, fmt.Errorf("failed to get verification token: %w", err)
	}

	return &verificationToken, nil
}

func (r *UserRepository) MarkTokenAsUsed(tokenID string) error {
	query := `UPDATE verification_tokens SET used = true WHERE id = $1`
	_, err := r.db.Exec(query, tokenID)
	return err
}