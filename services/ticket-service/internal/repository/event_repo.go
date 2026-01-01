package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/ticket-service/internal/models"
	"github.com/google/uuid"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event *models.Event) error {
	event.ID = uuid.New().String()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	formFieldsJSON, _ := json.Marshal(event.FormFields)

	query := `
		INSERT INTO events (id, creator_id, title, description, location, cover_image, 
			start_date, end_date, sale_start_date, sale_end_date, form_fields, 
			qr_code, event_code, status, currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := r.db.Exec(query,
		event.ID, event.CreatorID, event.Title, event.Description, event.Location,
		event.CoverImage, event.StartDate, event.EndDate, event.SaleStartDate,
		event.SaleEndDate, formFieldsJSON, event.QRCode, event.EventCode,
		event.Status, event.Currency, event.CreatedAt, event.UpdatedAt,
	)

	return err
}

func (r *EventRepository) GetByID(id string) (*models.Event, error) {
	query := `
		SELECT id, creator_id, title, description, location, cover_image,
			start_date, end_date, sale_start_date, sale_end_date, form_fields,
			qr_code, event_code, status, currency, created_at, updated_at
		FROM events WHERE id = $1
	`

	event := &models.Event{}
	var formFieldsJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&event.ID, &event.CreatorID, &event.Title, &event.Description,
		&event.Location, &event.CoverImage, &event.StartDate, &event.EndDate,
		&event.SaleStartDate, &event.SaleEndDate, &formFieldsJSON,
		&event.QRCode, &event.EventCode, &event.Status, &event.Currency,
		&event.CreatedAt, &event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(formFieldsJSON, &event.FormFields)
	return event, nil
}

func (r *EventRepository) GetByCode(code string) (*models.Event, error) {
	query := `
		SELECT id, creator_id, title, description, location, cover_image,
			start_date, end_date, sale_start_date, sale_end_date, form_fields,
			qr_code, event_code, status, currency, created_at, updated_at
		FROM events WHERE event_code = $1
	`

	event := &models.Event{}
	var formFieldsJSON []byte

	err := r.db.QueryRow(query, code).Scan(
		&event.ID, &event.CreatorID, &event.Title, &event.Description,
		&event.Location, &event.CoverImage, &event.StartDate, &event.EndDate,
		&event.SaleStartDate, &event.SaleEndDate, &formFieldsJSON,
		&event.QRCode, &event.EventCode, &event.Status, &event.Currency,
		&event.CreatedAt, &event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(formFieldsJSON, &event.FormFields)
	return event, nil
}

func (r *EventRepository) GetByCreator(creatorID string, limit, offset int) ([]*models.Event, error) {
	query := `
		SELECT id, creator_id, title, description, location, cover_image,
			start_date, end_date, sale_start_date, sale_end_date, form_fields,
			qr_code, event_code, status, currency, created_at, updated_at
		FROM events 
		WHERE creator_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, creatorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		event := &models.Event{}
		var formFieldsJSON []byte

		err := rows.Scan(
			&event.ID, &event.CreatorID, &event.Title, &event.Description,
			&event.Location, &event.CoverImage, &event.StartDate, &event.EndDate,
			&event.SaleStartDate, &event.SaleEndDate, &formFieldsJSON,
			&event.QRCode, &event.EventCode, &event.Status, &event.Currency,
			&event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(formFieldsJSON, &event.FormFields)
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetActive(limit, offset int) ([]*models.Event, error) {
	query := `
		SELECT id, creator_id, title, description, location, cover_image,
			start_date, end_date, sale_start_date, sale_end_date, form_fields,
			qr_code, event_code, status, currency, created_at, updated_at
		FROM events 
		WHERE status = 'active' AND sale_end_date > NOW()
		ORDER BY start_date ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		event := &models.Event{}
		var formFieldsJSON []byte

		err := rows.Scan(
			&event.ID, &event.CreatorID, &event.Title, &event.Description,
			&event.Location, &event.CoverImage, &event.StartDate, &event.EndDate,
			&event.SaleStartDate, &event.SaleEndDate, &formFieldsJSON,
			&event.QRCode, &event.EventCode, &event.Status, &event.Currency,
			&event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(formFieldsJSON, &event.FormFields)
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) Update(event *models.Event) error {
	event.UpdatedAt = time.Now()
	formFieldsJSON, _ := json.Marshal(event.FormFields)

	query := `
		UPDATE events SET
			title = $2, description = $3, location = $4, cover_image = $5,
			start_date = $6, end_date = $7, sale_start_date = $8, sale_end_date = $9,
			form_fields = $10, status = $11, updated_at = $12
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		event.ID, event.Title, event.Description, event.Location, event.CoverImage,
		event.StartDate, event.EndDate, event.SaleStartDate, event.SaleEndDate,
		formFieldsJSON, event.Status, event.UpdatedAt,
	)

	return err
}

func (r *EventRepository) UpdateStatus(id, status string) error {
	query := `UPDATE events SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id, status)
	return err
}

func (r *EventRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM events WHERE id = $1", id)
	return err
}
