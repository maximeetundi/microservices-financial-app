package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/models"
	"github.com/google/uuid"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	ticket.ID = uuid.New().String()
	ticket.CreatedAt = time.Now()

	formDataJSON, _ := json.Marshal(ticket.FormData)

	query := `
		INSERT INTO tickets (id, event_id, buyer_id, tier_id, tier_name, tier_icon,
			price, currency, form_data, qr_code, ticket_code, status, 
			transaction_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Exec(query,
		ticket.ID, ticket.EventID, ticket.BuyerID, ticket.TierID,
		ticket.TierName, ticket.TierIcon, ticket.Price, ticket.Currency,
		formDataJSON, ticket.QRCode, ticket.TicketCode, ticket.Status,
		ticket.TransactionID, ticket.CreatedAt,
	)

	return err
}

func (r *TicketRepository) GetByID(id string) (*models.Ticket, error) {
	query := `
		SELECT t.id, t.event_id, t.buyer_id, t.tier_id, t.tier_name, t.tier_icon,
			t.price, t.currency, t.form_data, t.qr_code, t.ticket_code, t.status,
			t.transaction_id, t.used_at, t.created_at,
			e.title, e.start_date, e.location
		FROM tickets t
		JOIN events e ON t.event_id = e.id
		WHERE t.id = $1
	`

	ticket := &models.Ticket{}
	var formDataJSON []byte
	var usedAt sql.NullTime
	var eventDate sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&ticket.ID, &ticket.EventID, &ticket.BuyerID, &ticket.TierID,
		&ticket.TierName, &ticket.TierIcon, &ticket.Price, &ticket.Currency,
		&formDataJSON, &ticket.QRCode, &ticket.TicketCode, &ticket.Status,
		&ticket.TransactionID, &usedAt, &ticket.CreatedAt,
		&ticket.EventTitle, &eventDate, &ticket.EventLocation,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(formDataJSON, &ticket.FormData)
	if usedAt.Valid {
		ticket.UsedAt = &usedAt.Time
	}
	if eventDate.Valid {
		ticket.EventDate = &eventDate.Time
	}

	return ticket, nil
}

func (r *TicketRepository) GetByCode(code string) (*models.Ticket, error) {
	query := `
		SELECT t.id, t.event_id, t.buyer_id, t.tier_id, t.tier_name, t.tier_icon,
			t.price, t.currency, t.form_data, t.qr_code, t.ticket_code, t.status,
			t.transaction_id, t.used_at, t.created_at,
			e.title, e.start_date, e.location
		FROM tickets t
		JOIN events e ON t.event_id = e.id
		WHERE t.ticket_code = $1
	`

	ticket := &models.Ticket{}
	var formDataJSON []byte
	var usedAt sql.NullTime
	var eventDate sql.NullTime

	err := r.db.QueryRow(query, code).Scan(
		&ticket.ID, &ticket.EventID, &ticket.BuyerID, &ticket.TierID,
		&ticket.TierName, &ticket.TierIcon, &ticket.Price, &ticket.Currency,
		&formDataJSON, &ticket.QRCode, &ticket.TicketCode, &ticket.Status,
		&ticket.TransactionID, &usedAt, &ticket.CreatedAt,
		&ticket.EventTitle, &eventDate, &ticket.EventLocation,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(formDataJSON, &ticket.FormData)
	if usedAt.Valid {
		ticket.UsedAt = &usedAt.Time
	}
	if eventDate.Valid {
		ticket.EventDate = &eventDate.Time
	}

	return ticket, nil
}

func (r *TicketRepository) GetByBuyer(buyerID string, limit, offset int) ([]*models.Ticket, error) {
	query := `
		SELECT t.id, t.event_id, t.buyer_id, t.tier_id, t.tier_name, t.tier_icon,
			t.price, t.currency, t.form_data, t.qr_code, t.ticket_code, t.status,
			t.transaction_id, t.used_at, t.created_at,
			e.title, e.start_date, e.location
		FROM tickets t
		JOIN events e ON t.event_id = e.id
		WHERE t.buyer_id = $1 AND t.status = 'paid'
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, buyerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*models.Ticket
	for rows.Next() {
		ticket := &models.Ticket{}
		var formDataJSON []byte
		var usedAt sql.NullTime
		var eventDate sql.NullTime

		err := rows.Scan(
			&ticket.ID, &ticket.EventID, &ticket.BuyerID, &ticket.TierID,
			&ticket.TierName, &ticket.TierIcon, &ticket.Price, &ticket.Currency,
			&formDataJSON, &ticket.QRCode, &ticket.TicketCode, &ticket.Status,
			&ticket.TransactionID, &usedAt, &ticket.CreatedAt,
			&ticket.EventTitle, &eventDate, &ticket.EventLocation,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(formDataJSON, &ticket.FormData)
		if usedAt.Valid {
			ticket.UsedAt = &usedAt.Time
		}
		if eventDate.Valid {
			ticket.EventDate = &eventDate.Time
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (r *TicketRepository) GetByEvent(eventID string, limit, offset int) ([]*models.Ticket, error) {
	query := `
		SELECT t.id, t.event_id, t.buyer_id, t.tier_id, t.tier_name, t.tier_icon,
			t.price, t.currency, t.form_data, t.qr_code, t.ticket_code, t.status,
			t.transaction_id, t.used_at, t.created_at
		FROM tickets t
		WHERE t.event_id = $1 AND t.status IN ('paid', 'used')
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, eventID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*models.Ticket
	for rows.Next() {
		ticket := &models.Ticket{}
		var formDataJSON []byte
		var usedAt sql.NullTime

		err := rows.Scan(
			&ticket.ID, &ticket.EventID, &ticket.BuyerID, &ticket.TierID,
			&ticket.TierName, &ticket.TierIcon, &ticket.Price, &ticket.Currency,
			&formDataJSON, &ticket.QRCode, &ticket.TicketCode, &ticket.Status,
			&ticket.TransactionID, &usedAt, &ticket.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(formDataJSON, &ticket.FormData)
		if usedAt.Valid {
			ticket.UsedAt = &usedAt.Time
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (r *TicketRepository) UpdateStatus(id, status string) error {
	query := `UPDATE tickets SET status = $2 WHERE id = $1`
	_, err := r.db.Exec(query, id, status)
	return err
}

func (r *TicketRepository) MarkAsUsed(id string) error {
	now := time.Now()
	query := `UPDATE tickets SET status = 'used', used_at = $2 WHERE id = $1`
	_, err := r.db.Exec(query, id, now)
	return err
}

func (r *TicketRepository) GetEventStats(eventID string) (*models.TicketStats, error) {
	// Get overall stats
	statsQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'paid' OR status = 'used') as sold,
			COUNT(*) FILTER (WHERE status = 'used') as used,
			COALESCE(SUM(price) FILTER (WHERE status = 'paid' OR status = 'used'), 0) as revenue
		FROM tickets WHERE event_id = $1
	`

	var stats models.TicketStats
	err := r.db.QueryRow(statsQuery, eventID).Scan(
		&stats.TotalTickets, &stats.SoldTickets, &stats.UsedTickets, &stats.TotalRevenue,
	)
	if err != nil {
		return nil, err
	}

	// Get tier breakdown
	tierQuery := `
		SELECT 
			tt.id, tt.name, tt.icon, tt.price, tt.quantity, tt.sold,
			COUNT(t.id) FILTER (WHERE t.status = 'used') as used,
			COALESCE(SUM(t.price) FILTER (WHERE t.status IN ('paid', 'used')), 0) as revenue
		FROM ticket_tiers tt
		LEFT JOIN tickets t ON t.tier_id = tt.id
		WHERE tt.event_id = $1
		GROUP BY tt.id, tt.name, tt.icon, tt.price, tt.quantity, tt.sold
		ORDER BY tt.sort_order ASC
	`

	rows, err := r.db.Query(tierQuery, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tierStats models.TierStats
		err := rows.Scan(
			&tierStats.TierID, &tierStats.TierName, &tierStats.TierIcon,
			&tierStats.Price, &tierStats.Quantity, &tierStats.Sold,
			&tierStats.Used, &tierStats.Revenue,
		)
		if err != nil {
			return nil, err
		}
		stats.TierBreakdown = append(stats.TierBreakdown, tierStats)
	}

	return &stats, nil
}
