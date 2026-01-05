package repository

import (
	"database/sql"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/models"
)

// ChatRepository handles chat message operations
type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(msg *models.ChatMessage) error {
	query := `
		INSERT INTO association_messages (association_id, sender_id, content, is_admin_only)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	return r.db.QueryRow(query, msg.AssociationID, msg.SenderID, msg.Content, msg.IsAdminOnly).
		Scan(&msg.ID, &msg.CreatedAt)
}

func (r *ChatRepository) GetByAssociation(associationID string, limit, offset int, includeAdminOnly bool) ([]*models.ChatMessage, error) {
	query := `
		SELECT m.id, m.association_id, m.sender_id, m.content, m.is_admin_only, m.created_at
		FROM association_messages m
		WHERE m.association_id = $1`
	if !includeAdminOnly {
		query += ` AND m.is_admin_only = false`
	}
	query += ` ORDER BY m.created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, associationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		msg := &models.ChatMessage{}
		if err := rows.Scan(&msg.ID, &msg.AssociationID, &msg.SenderID, &msg.Content, &msg.IsAdminOnly, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// SolidarityRepository handles solidarity event operations
type SolidarityRepository struct {
	db *sql.DB
}

func NewSolidarityRepository(db *sql.DB) *SolidarityRepository {
	return &SolidarityRepository{db: db}
}

func (r *SolidarityRepository) CreateEvent(event *models.SolidarityEvent) error {
	query := `
		INSERT INTO solidarity_events (association_id, event_type, beneficiary_id, title, description, target_amount, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`
	return r.db.QueryRow(query,
		event.AssociationID, event.EventType, event.BeneficiaryID, event.Title,
		event.Description, event.TargetAmount, event.Status, event.CreatedBy,
	).Scan(&event.ID, &event.CreatedAt)
}

func (r *SolidarityRepository) GetEvent(id string) (*models.SolidarityEvent, error) {
	query := `SELECT id, association_id, event_type, beneficiary_id, title, description, target_amount, collected_amount, status, created_by, created_at, closed_at
			  FROM solidarity_events WHERE id = $1`
	event := &models.SolidarityEvent{}
	var closedAt sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&event.ID, &event.AssociationID, &event.EventType, &event.BeneficiaryID, &event.Title,
		&event.Description, &event.TargetAmount, &event.CollectedAmount, &event.Status,
		&event.CreatedBy, &event.CreatedAt, &closedAt,
	)
	if closedAt.Valid {
		event.ClosedAt = &closedAt.Time
	}
	return event, err
}

func (r *SolidarityRepository) GetByAssociation(associationID string, status string) ([]*models.SolidarityEvent, error) {
	query := `SELECT id, association_id, event_type, beneficiary_id, title, description, target_amount, collected_amount, status, created_by, created_at
			  FROM solidarity_events WHERE association_id = $1`
	if status != "" {
		query += ` AND status = '` + status + `'`
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query, associationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.SolidarityEvent
	for rows.Next() {
		event := &models.SolidarityEvent{}
		if err := rows.Scan(&event.ID, &event.AssociationID, &event.EventType, &event.BeneficiaryID,
			&event.Title, &event.Description, &event.TargetAmount, &event.CollectedAmount,
			&event.Status, &event.CreatedBy, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *SolidarityRepository) UpdateCollectedAmount(id string, amount float64) error {
	_, err := r.db.Exec(`UPDATE solidarity_events SET collected_amount = collected_amount + $2 WHERE id = $1`, id, amount)
	return err
}

func (r *SolidarityRepository) CloseEvent(id string) error {
	_, err := r.db.Exec(`UPDATE solidarity_events SET status = 'closed', closed_at = $2 WHERE id = $1`, id, time.Now())
	return err
}

func (r *SolidarityRepository) AddContribution(contrib *models.SolidarityContribution) error {
	query := `
		INSERT INTO solidarity_contributions (event_id, contributor_id, amount, paid)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	return r.db.QueryRow(query, contrib.EventID, contrib.ContributorID, contrib.Amount, contrib.Paid).
		Scan(&contrib.ID, &contrib.CreatedAt)
}

func (r *SolidarityRepository) GetContributions(eventID string) ([]*models.SolidarityContribution, error) {
	query := `SELECT id, event_id, contributor_id, amount, paid, paid_at, created_at FROM solidarity_contributions WHERE event_id = $1`
	rows, err := r.db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contribs []*models.SolidarityContribution
	for rows.Next() {
		c := &models.SolidarityContribution{}
		var paidAt sql.NullTime
		if err := rows.Scan(&c.ID, &c.EventID, &c.ContributorID, &c.Amount, &c.Paid, &paidAt, &c.CreatedAt); err != nil {
			return nil, err
		}
		if paidAt.Valid {
			c.PaidAt = &paidAt.Time
		}
		contribs = append(contribs, c)
	}
	return contribs, nil
}

// CalledRoundRepository handles called tontine operations
type CalledRoundRepository struct {
	db *sql.DB
}

func NewCalledRoundRepository(db *sql.DB) *CalledRoundRepository {
	return &CalledRoundRepository{db: db}
}

func (r *CalledRoundRepository) Create(round *models.CalledRound) error {
	query := `
		INSERT INTO called_rounds (association_id, beneficiary_id, round_number, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	return r.db.QueryRow(query, round.AssociationID, round.BeneficiaryID, round.RoundNumber, round.Status).
		Scan(&round.ID, &round.CreatedAt)
}

func (r *CalledRoundRepository) GetByAssociation(associationID string, status string) ([]*models.CalledRound, error) {
	query := `SELECT id, association_id, beneficiary_id, round_number, total_collected, status, created_at
			  FROM called_rounds WHERE association_id = $1`
	if status != "" {
		query += ` AND status = '` + status + `'`
	}
	query += ` ORDER BY round_number DESC`

	rows, err := r.db.Query(query, associationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rounds []*models.CalledRound
	for rows.Next() {
		round := &models.CalledRound{}
		if err := rows.Scan(&round.ID, &round.AssociationID, &round.BeneficiaryID,
			&round.RoundNumber, &round.TotalCollected, &round.Status, &round.CreatedAt); err != nil {
			return nil, err
		}
		rounds = append(rounds, round)
	}
	return rounds, nil
}

func (r *CalledRoundRepository) GetNextRoundNumber(associationID string) (int, error) {
	var maxRound sql.NullInt64
	err := r.db.QueryRow(`SELECT MAX(round_number) FROM called_rounds WHERE association_id = $1`, associationID).Scan(&maxRound)
	if err != nil {
		return 1, err
	}
	if maxRound.Valid {
		return int(maxRound.Int64) + 1, nil
	}
	return 1, nil
}

func (r *CalledRoundRepository) AddPledge(pledge *models.CalledPledge) error {
	query := `
		INSERT INTO called_pledges (round_id, contributor_id, amount)
		VALUES ($1, $2, $3)
		ON CONFLICT (round_id, contributor_id) DO UPDATE SET amount = $3
		RETURNING id`
	return r.db.QueryRow(query, pledge.RoundID, pledge.ContributorID, pledge.Amount).Scan(&pledge.ID)
}

func (r *CalledRoundRepository) GetPledges(roundID string) ([]*models.CalledPledge, error) {
	query := `SELECT id, round_id, contributor_id, amount, paid, paid_at FROM called_pledges WHERE round_id = $1`
	rows, err := r.db.Query(query, roundID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pledges []*models.CalledPledge
	for rows.Next() {
		p := &models.CalledPledge{}
		var paidAt sql.NullTime
		if err := rows.Scan(&p.ID, &p.RoundID, &p.ContributorID, &p.Amount, &p.Paid, &paidAt); err != nil {
			return nil, err
		}
		if paidAt.Valid {
			p.PaidAt = &paidAt.Time
		}
		pledges = append(pledges, p)
	}
	return pledges, nil
}

func (r *CalledRoundRepository) MarkPledgePaid(pledgeID string) error {
	_, err := r.db.Exec(`UPDATE called_pledges SET paid = true, paid_at = $2 WHERE id = $1`, pledgeID, time.Now())
	return err
}

func (r *CalledRoundRepository) UpdateTotalCollected(roundID string, amount float64) error {
	_, err := r.db.Exec(`UPDATE called_rounds SET total_collected = total_collected + $2 WHERE id = $1`, roundID, amount)
	return err
}

func (r *CalledRoundRepository) CloseRound(roundID string) error {
	_, err := r.db.Exec(`UPDATE called_rounds SET status = 'closed', closed_at = $2 WHERE id = $1`, roundID, time.Now())
	return err
}
