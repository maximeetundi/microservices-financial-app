package repository

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/association-service/internal/models"
)

// RoleRepository handles custom role operations
type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *models.AssociationRole) error {
	permJSON, _ := json.Marshal(role.Permissions)
	query := `
		INSERT INTO association_roles (association_id, name, permissions, is_default)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	return r.db.QueryRow(query, role.AssociationID, role.Name, permJSON, role.IsDefault).
		Scan(&role.ID, &role.CreatedAt)
}

func (r *RoleRepository) GetByAssociation(associationID string) ([]*models.AssociationRole, error) {
	query := `SELECT id, association_id, name, permissions, is_default, created_at 
			  FROM association_roles WHERE association_id = $1 ORDER BY is_default DESC, name`
	rows, err := r.db.Query(query, associationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.AssociationRole
	for rows.Next() {
		role := &models.AssociationRole{}
		var permJSON []byte
		if err := rows.Scan(&role.ID, &role.AssociationID, &role.Name, &permJSON, &role.IsDefault, &role.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(permJSON, &role.Permissions)
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM association_roles WHERE id = $1`, id)
	return err
}

func (r *RoleRepository) Update(role *models.AssociationRole) error {
	permJSON, _ := json.Marshal(role.Permissions)
	_, err := r.db.Exec(`UPDATE association_roles SET name = $2, permissions = $3 WHERE id = $1`,
		role.ID, role.Name, permJSON)
	return err
}

// ApprovalRepository handles multi-signature approval operations
type ApprovalRepository struct {
	db *sql.DB
}

func NewApprovalRepository(db *sql.DB) *ApprovalRepository {
	return &ApprovalRepository{db: db}
}

// SetApprovers sets the 5 designated approvers for an association
func (r *ApprovalRepository) SetApprovers(associationID string, memberIDs []string) error {
	// Delete existing approvers
	_, _ = r.db.Exec(`DELETE FROM association_approvers WHERE association_id = $1`, associationID)

	// Insert new approvers
	for i, memberID := range memberIDs {
		_, err := r.db.Exec(`
			INSERT INTO association_approvers (association_id, member_id, position)
			VALUES ($1, $2, $3)`, associationID, memberID, i+1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ApprovalRepository) GetApprovers(associationID string) ([]*models.Approver, error) {
	query := `
		SELECT a.id, a.association_id, a.member_id, a.position, a.created_at
		FROM association_approvers a
		WHERE a.association_id = $1
		ORDER BY a.position`
	rows, err := r.db.Query(query, associationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var approvers []*models.Approver
	for rows.Next() {
		a := &models.Approver{}
		if err := rows.Scan(&a.ID, &a.AssociationID, &a.MemberID, &a.Position, &a.CreatedAt); err != nil {
			return nil, err
		}
		approvers = append(approvers, a)
	}
	return approvers, nil
}

func (r *ApprovalRepository) IsApprover(associationID, memberID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM association_approvers WHERE association_id = $1 AND member_id = $2`,
		associationID, memberID).Scan(&count)
	return count > 0, err
}

// CreateRequest creates a new approval request
func (r *ApprovalRepository) CreateRequest(req *models.ApprovalRequest) error {
	query := `
		INSERT INTO approval_requests (association_id, request_type, reference_id, amount, requester_id, status, required_approvals, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query,
		req.AssociationID, req.RequestType, req.ReferenceID, req.Amount, req.RequesterID,
		req.Status, req.RequiredApprovals, req.Description,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *ApprovalRepository) GetRequest(id string) (*models.ApprovalRequest, error) {
	query := `SELECT id, association_id, request_type, reference_id, amount, requester_id, status, required_approvals, description, created_at, updated_at
			  FROM approval_requests WHERE id = $1`
	req := &models.ApprovalRequest{}
	err := r.db.QueryRow(query, id).Scan(
		&req.ID, &req.AssociationID, &req.RequestType, &req.ReferenceID, &req.Amount,
		&req.RequesterID, &req.Status, &req.RequiredApprovals, &req.Description, &req.CreatedAt, &req.UpdatedAt,
	)
	return req, err
}

func (r *ApprovalRepository) GetPendingRequests(associationID string) ([]*models.ApprovalRequest, error) {
	query := `SELECT id, association_id, request_type, reference_id, amount, requester_id, status, required_approvals, description, created_at, updated_at
			  FROM approval_requests WHERE association_id = $1 AND status = 'pending' ORDER BY created_at DESC`
	rows, err := r.db.Query(query, associationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.ApprovalRequest
	for rows.Next() {
		req := &models.ApprovalRequest{}
		if err := rows.Scan(&req.ID, &req.AssociationID, &req.RequestType, &req.ReferenceID, &req.Amount,
			&req.RequesterID, &req.Status, &req.RequiredApprovals, &req.Description, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *ApprovalRepository) UpdateRequestStatus(id string, status models.ApprovalRequestStatus) error {
	_, err := r.db.Exec(`UPDATE approval_requests SET status = $2, updated_at = $3 WHERE id = $1`,
		id, status, time.Now())
	return err
}

// Vote records a vote on an approval request
func (r *ApprovalRepository) Vote(vote *models.ApprovalVote) error {
	query := `
		INSERT INTO approval_votes (request_id, approver_id, vote, comment)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (request_id, approver_id) DO UPDATE SET vote = $3, comment = $4, voted_at = CURRENT_TIMESTAMP
		RETURNING id, voted_at`
	return r.db.QueryRow(query, vote.RequestID, vote.ApproverID, vote.Vote, vote.Comment).
		Scan(&vote.ID, &vote.VotedAt)
}

func (r *ApprovalRepository) GetVotes(requestID string) ([]*models.ApprovalVote, error) {
	query := `SELECT id, request_id, approver_id, vote, comment, voted_at FROM approval_votes WHERE request_id = $1`
	rows, err := r.db.Query(query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []*models.ApprovalVote
	for rows.Next() {
		v := &models.ApprovalVote{}
		if err := rows.Scan(&v.ID, &v.RequestID, &v.ApproverID, &v.Vote, &v.Comment, &v.VotedAt); err != nil {
			return nil, err
		}
		votes = append(votes, v)
	}
	return votes, nil
}

func (r *ApprovalRepository) CountVotes(requestID string) (approvals int, rejections int, err error) {
	err = r.db.QueryRow(`SELECT COUNT(*) FROM approval_votes WHERE request_id = $1 AND vote = 'approve'`, requestID).Scan(&approvals)
	if err != nil {
		return
	}
	err = r.db.QueryRow(`SELECT COUNT(*) FROM approval_votes WHERE request_id = $1 AND vote = 'reject'`, requestID).Scan(&rejections)
	return
}

func (r *ApprovalRepository) HasVoted(requestID, memberID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM approval_votes WHERE request_id = $1 AND approver_id = $2`,
		requestID, memberID).Scan(&count)
	return count > 0, err
}
