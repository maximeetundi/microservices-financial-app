package repository

import (
	"database/sql"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/card-service/internal/models"
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(card *models.Card) error {
	query := `
		INSERT INTO cards (id, user_id, card_number, card_type, currency, balance, daily_limit, monthly_limit, is_active, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(query,
		card.ID,
		card.UserID,
		card.CardNumber,
		card.CardType,
		card.Currency,
		card.Balance,
		card.DailyLimit,
		card.MonthlyLimit,
		card.IsActive,
		card.ExpiresAt,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *CardRepository) GetByID(id string) (*models.Card, error) {
	query := `
		SELECT id, user_id, card_number, card_type, currency, balance, daily_limit, monthly_limit, is_active, expires_at, created_at, updated_at
		FROM cards WHERE id = $1
	`
	var card models.Card
	err := r.db.QueryRow(query, id).Scan(
		&card.ID,
		&card.UserID,
		&card.CardNumber,
		&card.CardType,
		&card.Currency,
		&card.Balance,
		&card.DailyLimit,
		&card.MonthlyLimit,
		&card.IsActive,
		&card.ExpiresAt,
		&card.CreatedAt,
		&card.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *CardRepository) GetByUserID(userID string) ([]models.Card, error) {
	return r.GetUserCards(userID)
}

func (r *CardRepository) GetUserCards(userID string) ([]models.Card, error) {
	query := `
		SELECT id, user_id, card_number, card_type, currency, balance, daily_limit, monthly_limit, is_active, expires_at, created_at, updated_at
		FROM cards WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var c models.Card
		err := rows.Scan(
			&c.ID, &c.UserID, &c.CardNumber, &c.CardType, &c.Currency,
			&c.Balance, &c.DailyLimit, &c.MonthlyLimit, &c.IsActive, &c.ExpiresAt,
			&c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func (r *CardRepository) Update(card *models.Card) error {
	query := `
		UPDATE cards SET balance = $1, daily_limit = $2, monthly_limit = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(query, card.Balance, card.DailyLimit, card.MonthlyLimit, card.IsActive, time.Now(), card.ID)
	return err
}

func (r *CardRepository) UpdateBalance(cardID string, amount float64) error {
	query := `UPDATE cards SET balance = balance + $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, amount, time.Now(), cardID)
	return err
}

func (r *CardRepository) Delete(id string) error {
	query := `DELETE FROM cards WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Gift card methods
func (r *CardRepository) CreateGiftCard(giftCard *models.GiftCard) error {
	query := `
		INSERT INTO gift_cards (id, code, sender_id, recipient_email, recipient_phone, amount, currency, message, design, status, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(query,
		giftCard.ID,
		giftCard.Code,
		giftCard.SenderID,
		giftCard.RecipientEmail,
		giftCard.RecipientPhone,
		giftCard.Amount,
		giftCard.Currency,
		giftCard.Message,
		giftCard.Design,
		giftCard.Status,
		giftCard.ExpiresAt,
		time.Now(),
	)
	return err
}

func (r *CardRepository) GetGiftCardByCode(code string) (*models.GiftCard, error) {
	query := `
		SELECT id, code, sender_id, recipient_email, recipient_phone, amount, currency, message, design, status, redeemed_by, redeemed_at, expires_at, created_at
		FROM gift_cards WHERE code = $1
	`
	var gc models.GiftCard
	err := r.db.QueryRow(query, code).Scan(
		&gc.ID,
		&gc.Code,
		&gc.SenderID,
		&gc.RecipientEmail,
		&gc.RecipientPhone,
		&gc.Amount,
		&gc.Currency,
		&gc.Message,
		&gc.Design,
		&gc.Status,
		&gc.RedeemedBy,
		&gc.RedeemedAt,
		&gc.ExpiresAt,
		&gc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &gc, nil
}

func (r *CardRepository) UpdateGiftCard(giftCard *models.GiftCard) error {
	query := `
		UPDATE gift_cards SET status = $1, redeemed_by = $2, redeemed_at = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(query, giftCard.Status, giftCard.RedeemedBy, giftCard.RedeemedAt, giftCard.ID)
	return err
}

type CardTransactionRepository struct {
	db *sql.DB
}

func NewCardTransactionRepository(db *sql.DB) *CardTransactionRepository {
	return &CardTransactionRepository{db: db}
}

func (r *CardTransactionRepository) Create(tx *models.CardTransaction) error {
	query := `
		INSERT INTO card_transactions (id, card_id, user_id, transaction_type, amount, currency, fee, status, merchant_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(query,
		tx.ID,
		tx.CardID,
		tx.UserID,
		tx.TransactionType,
		tx.Amount,
		tx.Currency,
		tx.Fee,
		tx.Status,
		tx.MerchantName,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *CardTransactionRepository) UpdateStatus(transactionID, status string) error {
	query := `UPDATE card_transactions SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), transactionID)
	return err
}

func (r *CardTransactionRepository) GetByCardID(cardID string, limit, offset int) ([]models.CardTransaction, error) {
	query := `
		SELECT id, card_id, user_id, transaction_type, amount, currency, fee, status, merchant_name, created_at
		FROM card_transactions 
		WHERE card_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, cardID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.CardTransaction
	for rows.Next() {
		var t models.CardTransaction
		err := rows.Scan(&t.ID, &t.CardID, &t.UserID, &t.TransactionType, &t.Amount, &t.Currency, &t.Fee, &t.Status, &t.MerchantName, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
