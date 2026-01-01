package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/models"
	"github.com/google/uuid"
)

type TierRepository struct {
	db *sql.DB
}

func NewTierRepository(db *sql.DB) *TierRepository {
	return &TierRepository{db: db}
}

func (r *TierRepository) Create(tier *models.TicketTier) error {
	tier.ID = uuid.New().String()
	tier.CreatedAt = time.Now()

	benefitsJSON, _ := json.Marshal(tier.Benefits)

	query := `
		INSERT INTO ticket_tiers (id, event_id, name, icon, price, quantity, sold, 
			description, benefits, color, sort_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.Exec(query,
		tier.ID, tier.EventID, tier.Name, tier.Icon, tier.Price,
		tier.Quantity, tier.Sold, tier.Description, benefitsJSON,
		tier.Color, tier.SortOrder, tier.CreatedAt,
	)

	return err
}

func (r *TierRepository) GetByEventID(eventID string) ([]*models.TicketTier, error) {
	query := `
		SELECT id, event_id, name, icon, price, quantity, sold, 
			description, benefits, color, sort_order, created_at
		FROM ticket_tiers 
		WHERE event_id = $1
		ORDER BY sort_order ASC, price ASC
	`

	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tiers []*models.TicketTier
	for rows.Next() {
		tier := &models.TicketTier{}
		var benefitsJSON []byte

		err := rows.Scan(
			&tier.ID, &tier.EventID, &tier.Name, &tier.Icon, &tier.Price,
			&tier.Quantity, &tier.Sold, &tier.Description, &benefitsJSON,
			&tier.Color, &tier.SortOrder, &tier.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(benefitsJSON, &tier.Benefits)
		tiers = append(tiers, tier)
	}

	return tiers, nil
}

func (r *TierRepository) GetByID(id string) (*models.TicketTier, error) {
	query := `
		SELECT id, event_id, name, icon, price, quantity, sold, 
			description, benefits, color, sort_order, created_at
		FROM ticket_tiers WHERE id = $1
	`

	tier := &models.TicketTier{}
	var benefitsJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&tier.ID, &tier.EventID, &tier.Name, &tier.Icon, &tier.Price,
		&tier.Quantity, &tier.Sold, &tier.Description, &benefitsJSON,
		&tier.Color, &tier.SortOrder, &tier.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(benefitsJSON, &tier.Benefits)
	return tier, nil
}

func (r *TierRepository) IncrementSold(tierID string) error {
	query := `UPDATE ticket_tiers SET sold = sold + 1 WHERE id = $1`
	_, err := r.db.Exec(query, tierID)
	return err
}

func (r *TierRepository) DecrementSold(tierID string) error {
	query := `UPDATE ticket_tiers SET sold = sold - 1 WHERE id = $1 AND sold > 0`
	_, err := r.db.Exec(query, tierID)
	return err
}

func (r *TierRepository) DeleteByEventID(eventID string) error {
	_, err := r.db.Exec("DELETE FROM ticket_tiers WHERE event_id = $1", eventID)
	return err
}

func (r *TierRepository) CheckAvailability(tierID string) (bool, error) {
	query := `SELECT quantity, sold FROM ticket_tiers WHERE id = $1`
	
	var quantity, sold int
	err := r.db.QueryRow(query, tierID).Scan(&quantity, &sold)
	if err != nil {
		return false, err
	}

	// -1 means unlimited
	if quantity == -1 {
		return true, nil
	}

	return sold < quantity, nil
}
