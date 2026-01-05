package repository

import (
	"database/sql"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/models"
)

type AssociationRepository struct {
	db *sql.DB
}

func NewAssociationRepository(db *sql.DB) *AssociationRepository {
	return &AssociationRepository{db: db}
}

func (r *AssociationRepository) Create(a *models.Association) error {
	query := `
		INSERT INTO associations (name, type, description, rules, currency, status, created_by, creator_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		a.Name, a.Type, a.Description, a.Rules, a.Currency, a.Status, a.CreatedBy, a.CreatedBy,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *AssociationRepository) GetByID(id string) (*models.Association, error) {
	query := `SELECT id, name, type, description, rules, total_members, treasury_balance, 
			  currency, status, created_by, created_at, updated_at 
			  FROM associations WHERE id = $1`

	a := &models.Association{}
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.Name, &a.Type, &a.Description, &a.Rules, &a.TotalMembers,
		&a.TreasuryBalance, &a.Currency, &a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *AssociationRepository) GetByUser(userID string, limit, offset int) ([]*models.Association, error) {
	query := `
		SELECT DISTINCT a.id, a.name, a.type, a.description, a.rules, a.total_members, 
			   a.treasury_balance, a.currency, a.status, a.created_by, a.created_at, a.updated_at
		FROM associations a
		LEFT JOIN members m ON a.id = m.association_id
		WHERE a.created_by = $1 OR (m.user_id = $1 AND m.status = 'active')
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var associations []*models.Association
	for rows.Next() {
		a := &models.Association{}
		if err := rows.Scan(
			&a.ID, &a.Name, &a.Type, &a.Description, &a.Rules, &a.TotalMembers,
			&a.TreasuryBalance, &a.Currency, &a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		associations = append(associations, a)
	}
	return associations, nil
}

func (r *AssociationRepository) Update(a *models.Association) error {
	query := `
		UPDATE associations 
		SET name = $2, description = $3, rules = $4, status = $5, updated_at = $6
		WHERE id = $1`

	_, err := r.db.Exec(query, a.ID, a.Name, a.Description, a.Rules, a.Status, time.Now())
	return err
}

func (r *AssociationRepository) UpdateTreasuryBalance(id string, amount float64) error {
	query := `UPDATE associations SET treasury_balance = treasury_balance + $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, amount, time.Now())
	return err
}

func (r *AssociationRepository) UpdateMemberCount(id string, delta int) error {
	query := `UPDATE associations SET total_members = total_members + $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, delta, time.Now())
	return err
}

func (r *AssociationRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM associations WHERE id = $1`, id)
	return err
}

// GetAll for admin
func (r *AssociationRepository) GetAll(limit, offset int) ([]*models.Association, error) {
	query := `
		SELECT id, name, type, description, rules, total_members, treasury_balance, 
			   currency, status, created_by, created_at, updated_at
		FROM associations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var associations []*models.Association
	for rows.Next() {
		a := &models.Association{}
		if err := rows.Scan(
			&a.ID, &a.Name, &a.Type, &a.Description, &a.Rules, &a.TotalMembers,
			&a.TreasuryBalance, &a.Currency, &a.Status, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		associations = append(associations, a)
	}
	return associations, nil
}

// === Member Repository ===

type MemberRepository struct {
	db *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) Create(m *models.Member) error {
	query := `
		INSERT INTO members (association_id, user_id, role, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, join_date, created_at, updated_at`

	return r.db.QueryRow(query, m.AssociationID, m.UserID, m.Role, m.Status).
		Scan(&m.ID, &m.JoinDate, &m.CreatedAt, &m.UpdatedAt)
}

func (r *MemberRepository) GetByID(id string) (*models.Member, error) {
	query := `SELECT id, association_id, user_id, role, join_date, status, 
			  contributions_paid, loans_received, created_at, updated_at 
			  FROM members WHERE id = $1`

	m := &models.Member{}
	err := r.db.QueryRow(query, id).Scan(
		&m.ID, &m.AssociationID, &m.UserID, &m.Role, &m.JoinDate, &m.Status,
		&m.ContributionsPaid, &m.LoansReceived, &m.CreatedAt, &m.UpdatedAt,
	)
	return m, err
}

func (r *MemberRepository) GetByAssociationAndUser(associationID, userID string) (*models.Member, error) {
	query := `SELECT id, association_id, user_id, role, join_date, status, 
			  contributions_paid, loans_received, created_at, updated_at 
			  FROM members WHERE association_id = $1 AND user_id = $2`

	m := &models.Member{}
	err := r.db.QueryRow(query, associationID, userID).Scan(
		&m.ID, &m.AssociationID, &m.UserID, &m.Role, &m.JoinDate, &m.Status,
		&m.ContributionsPaid, &m.LoansReceived, &m.CreatedAt, &m.UpdatedAt,
	)
	return m, err
}

func (r *MemberRepository) GetByAssociation(associationID string) ([]*models.Member, error) {
	query := `SELECT id, association_id, user_id, role, join_date, status, 
			  contributions_paid, loans_received, created_at, updated_at 
			  FROM members WHERE association_id = $1 ORDER BY role, join_date`

	rows, err := r.db.Query(query, associationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.Member
	for rows.Next() {
		m := &models.Member{}
		if err := rows.Scan(
			&m.ID, &m.AssociationID, &m.UserID, &m.Role, &m.JoinDate, &m.Status,
			&m.ContributionsPaid, &m.LoansReceived, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *MemberRepository) UpdateRole(id string, role models.MemberRole) error {
	_, err := r.db.Exec(`UPDATE members SET role = $2, updated_at = $3 WHERE id = $1`, id, role, time.Now())
	return err
}

func (r *MemberRepository) UpdateStatus(id string, status models.MemberStatus) error {
	_, err := r.db.Exec(`UPDATE members SET status = $2, updated_at = $3 WHERE id = $1`, id, status, time.Now())
	return err
}

func (r *MemberRepository) UpdateContributions(id string, amount float64) error {
	_, err := r.db.Exec(`UPDATE members SET contributions_paid = contributions_paid + $2, updated_at = $3 WHERE id = $1`, id, amount, time.Now())
	return err
}

func (r *MemberRepository) UpdateLoansReceived(id string, amount float64) error {
	_, err := r.db.Exec(`UPDATE members SET loans_received = loans_received + $2, updated_at = $3 WHERE id = $1`, id, amount, time.Now())
	return err
}

func (r *MemberRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM members WHERE id = $1`, id)
	return err
}

// === Treasury Repository ===

type TreasuryRepository struct {
	db *sql.DB
}

func NewTreasuryRepository(db *sql.DB) *TreasuryRepository {
	return &TreasuryRepository{db: db}
}

func (r *TreasuryRepository) CreateTransaction(t *models.TreasuryTransaction) error {
	query := `
		INSERT INTO treasury_transactions (association_id, type, amount, from_member_id, to_member_id, description, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	var fromID, toID interface{}
	if t.FromMemberID != "" {
		fromID = t.FromMemberID
	}
	if t.ToMemberID != "" {
		toID = t.ToMemberID
	}

	return r.db.QueryRow(query,
		t.AssociationID, t.Type, t.Amount, fromID, toID, t.Description, t.Status, t.CreatedBy,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TreasuryRepository) GetByAssociation(associationID string, limit, offset int) ([]*models.TreasuryTransaction, error) {
	query := `
		SELECT id, association_id, type, amount, from_member_id, to_member_id, description, status, created_by, created_at, updated_at
		FROM treasury_transactions 
		WHERE association_id = $1 
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, associationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.TreasuryTransaction
	for rows.Next() {
		t := &models.TreasuryTransaction{}
		var fromID, toID sql.NullString
		if err := rows.Scan(
			&t.ID, &t.AssociationID, &t.Type, &t.Amount, &fromID, &toID,
			&t.Description, &t.Status, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if fromID.Valid {
			t.FromMemberID = fromID.String
		}
		if toID.Valid {
			t.ToMemberID = toID.String
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *TreasuryRepository) GetReport(associationID string) (*models.TreasuryReport, error) {
	report := &models.TreasuryReport{}

	// Get totals by type
	query := `
		SELECT type, COALESCE(SUM(amount), 0) as total
		FROM treasury_transactions
		WHERE association_id = $1 AND status = 'completed'
		GROUP BY type`

	rows, err := r.db.Query(query, associationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var txType string
		var total float64
		if err := rows.Scan(&txType, &total); err != nil {
			return nil, err
		}
		switch txType {
		case "contribution":
			report.TotalContributions = total
		case "loan":
			report.TotalLoans = total
		case "repayment":
			report.TotalRepayments = total
		case "distribution":
			report.TotalDistributions = total
		case "expense":
			report.TotalExpenses = total
		}
	}

	// Get balance
	var balance sql.NullFloat64
	err = r.db.QueryRow(`SELECT treasury_balance FROM associations WHERE id = $1`, associationID).Scan(&balance)
	if err != nil {
		return nil, err
	}
	report.TotalBalance = balance.Float64

	// Get recent transactions
	report.Transactions, _ = r.GetByAssociation(associationID, 20, 0)

	return report, nil
}

func (r *TreasuryRepository) UpdateStatus(id string, status models.TransactionStatus) error {
	_, err := r.db.Exec(`UPDATE treasury_transactions SET status = $2, updated_at = $3 WHERE id = $1`, id, status, time.Now())
	return err
}

// === Meeting Repository ===

type MeetingRepository struct {
	db *sql.DB
}

func NewMeetingRepository(db *sql.DB) *MeetingRepository {
	return &MeetingRepository{db: db}
}

func (r *MeetingRepository) Create(m *models.Meeting) error {
	query := `
		INSERT INTO meetings (association_id, title, date, location, agenda, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		m.AssociationID, m.Title, m.Date, m.Location, m.Agenda, m.Status, m.CreatedBy,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *MeetingRepository) GetByID(id string) (*models.Meeting, error) {
	query := `SELECT id, association_id, title, date, location, agenda, minutes, attendance, status, created_by, created_at, updated_at
			  FROM meetings WHERE id = $1`

	m := &models.Meeting{}
	err := r.db.QueryRow(query, id).Scan(
		&m.ID, &m.AssociationID, &m.Title, &m.Date, &m.Location, &m.Agenda,
		&m.Minutes, &m.Attendance, &m.Status, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
	)
	return m, err
}

func (r *MeetingRepository) GetByAssociation(associationID string, limit, offset int) ([]*models.Meeting, error) {
	query := `
		SELECT id, association_id, title, date, location, agenda, minutes, attendance, status, created_by, created_at, updated_at
		FROM meetings WHERE association_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, associationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []*models.Meeting
	for rows.Next() {
		m := &models.Meeting{}
		if err := rows.Scan(
			&m.ID, &m.AssociationID, &m.Title, &m.Date, &m.Location, &m.Agenda,
			&m.Minutes, &m.Attendance, &m.Status, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		meetings = append(meetings, m)
	}
	return meetings, nil
}

func (r *MeetingRepository) UpdateAttendance(id string, attendance models.JSONB) error {
	_, err := r.db.Exec(`UPDATE meetings SET attendance = $2, updated_at = $3 WHERE id = $1`, id, attendance, time.Now())
	return err
}

func (r *MeetingRepository) UpdateMinutes(id, minutes string) error {
	_, err := r.db.Exec(`UPDATE meetings SET minutes = $2, updated_at = $3 WHERE id = $1`, id, minutes, time.Now())
	return err
}

func (r *MeetingRepository) UpdateStatus(id string, status models.MeetingStatus) error {
	_, err := r.db.Exec(`UPDATE meetings SET status = $2, updated_at = $3 WHERE id = $1`, id, status, time.Now())
	return err
}

// === Loan Repository ===

type LoanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) Create(l *models.Loan) error {
	query := `
		INSERT INTO loans (association_id, borrower_id, amount, interest_rate, duration, start_date, end_date, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		l.AssociationID, l.BorrowerID, l.Amount, l.InterestRate, l.Duration, l.StartDate, l.EndDate, l.Status,
	).Scan(&l.ID, &l.CreatedAt, &l.UpdatedAt)
}

func (r *LoanRepository) GetByID(id string) (*models.Loan, error) {
	query := `SELECT id, association_id, borrower_id, amount, interest_rate, duration, 
			  start_date, end_date, repayments, status, approved_by, created_at, updated_at
			  FROM loans WHERE id = $1`

	l := &models.Loan{}
	var approvedBy sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&l.ID, &l.AssociationID, &l.BorrowerID, &l.Amount, &l.InterestRate, &l.Duration,
		&l.StartDate, &l.EndDate, &l.Repayments, &l.Status, &approvedBy, &l.CreatedAt, &l.UpdatedAt,
	)
	if approvedBy.Valid {
		l.ApprovedBy = approvedBy.String
	}
	return l, err
}

func (r *LoanRepository) GetByAssociation(associationID string, limit, offset int) ([]*models.Loan, error) {
	query := `
		SELECT id, association_id, borrower_id, amount, interest_rate, duration, 
			   start_date, end_date, repayments, status, approved_by, created_at, updated_at
		FROM loans WHERE association_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, associationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*models.Loan
	for rows.Next() {
		l := &models.Loan{}
		var approvedBy sql.NullString
		if err := rows.Scan(
			&l.ID, &l.AssociationID, &l.BorrowerID, &l.Amount, &l.InterestRate, &l.Duration,
			&l.StartDate, &l.EndDate, &l.Repayments, &l.Status, &approvedBy, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if approvedBy.Valid {
			l.ApprovedBy = approvedBy.String
		}
		loans = append(loans, l)
	}
	return loans, nil
}

func (r *LoanRepository) GetByBorrower(borrowerID string) ([]*models.Loan, error) {
	query := `
		SELECT id, association_id, borrower_id, amount, interest_rate, duration, 
			   start_date, end_date, repayments, status, approved_by, created_at, updated_at
		FROM loans WHERE borrower_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, borrowerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*models.Loan
	for rows.Next() {
		l := &models.Loan{}
		var approvedBy sql.NullString
		if err := rows.Scan(
			&l.ID, &l.AssociationID, &l.BorrowerID, &l.Amount, &l.InterestRate, &l.Duration,
			&l.StartDate, &l.EndDate, &l.Repayments, &l.Status, &approvedBy, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if approvedBy.Valid {
			l.ApprovedBy = approvedBy.String
		}
		loans = append(loans, l)
	}
	return loans, nil
}

func (r *LoanRepository) Approve(id, approvedBy string, status models.LoanStatus) error {
	_, err := r.db.Exec(`UPDATE loans SET status = $2, approved_by = $3, updated_at = $4 WHERE id = $1`,
		id, status, approvedBy, time.Now())
	return err
}

func (r *LoanRepository) UpdateRepayments(id string, repayments models.JSONB) error {
	_, err := r.db.Exec(`UPDATE loans SET repayments = $2, updated_at = $3 WHERE id = $1`, id, repayments, time.Now())
	return err
}

func (r *LoanRepository) UpdateStatus(id string, status models.LoanStatus) error {
	_, err := r.db.Exec(`UPDATE loans SET status = $2, updated_at = $3 WHERE id = $1`, id, status, time.Now())
	return err
}
