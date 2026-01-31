package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/crypto-bank/microservices-financial-app/services/admin-service/internal/models"
)

type AdminRepository struct {
	adminDB *sql.DB
	mainDB  *sql.DB
}

func NewAdminRepository(adminDB, mainDB *sql.DB) *AdminRepository {
	return &AdminRepository{
		adminDB: adminDB,
		mainDB:  mainDB,
	}
}

// ========== Admin Users ==========

func (r *AdminRepository) CreateAdminUser(admin *models.AdminUser) error {
	query := `
		INSERT INTO admin_users (id, email, password_hash, first_name, last_name, role_id, is_active, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.adminDB.Exec(query,
		admin.ID, admin.Email, admin.PasswordHash,
		admin.FirstName, admin.LastName, admin.RoleID,
		admin.IsActive, admin.CreatedAt, admin.UpdatedAt, admin.CreatedBy,
	)
	return err
}

func (r *AdminRepository) GetAdminByEmail(email string) (*models.AdminUser, error) {
	admin := &models.AdminUser{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, role_id, is_active, last_login_at, created_at, updated_at
		FROM admin_users WHERE email = $1
	`
	err := r.adminDB.QueryRow(query, email).Scan(
		&admin.ID, &admin.Email, &admin.PasswordHash,
		&admin.FirstName, &admin.LastName, &admin.RoleID,
		&admin.IsActive, &admin.LastLoginAt, &admin.CreatedAt, &admin.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (r *AdminRepository) GetAdminByID(id string) (*models.AdminUser, error) {
	admin := &models.AdminUser{}
	query := `
		SELECT id, email, password_hash, first_name, last_name, role_id, is_active, last_login_at, created_at, updated_at
		FROM admin_users WHERE id = $1
	`
	err := r.adminDB.QueryRow(query, id).Scan(
		&admin.ID, &admin.Email, &admin.PasswordHash,
		&admin.FirstName, &admin.LastName, &admin.RoleID,
		&admin.IsActive, &admin.LastLoginAt, &admin.CreatedAt, &admin.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (r *AdminRepository) GetAllAdmins(limit, offset int) ([]models.AdminUser, error) {
	query := `
		SELECT a.id, a.email, a.first_name, a.last_name, a.role_id, r.name as role_name, a.is_active, a.last_login_at, a.created_at
		FROM admin_users a
		LEFT JOIN admin_roles r ON a.role_id = r.id
		ORDER BY a.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.adminDB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []models.AdminUser
	for rows.Next() {
		var admin models.AdminUser
		var roleName sql.NullString
		err := rows.Scan(
			&admin.ID, &admin.Email, &admin.FirstName, &admin.LastName,
			&admin.RoleID, &roleName, &admin.IsActive, &admin.LastLoginAt, &admin.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if roleName.Valid {
			admin.Role = &models.AdminRole{Name: roleName.String}
		}
		admins = append(admins, admin)
	}
	return admins, nil
}

func (r *AdminRepository) UpdateAdmin(id string, updates map[string]interface{}) error {
	query := "UPDATE admin_users SET updated_at = $1"
	args := []interface{}{time.Now()}
	argCount := 2

	for key, value := range updates {
		query += fmt.Sprintf(", %s = $%d", key, argCount)
		args = append(args, value)
		argCount++
	}
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, id)

	_, err := r.adminDB.Exec(query, args...)
	return err
}

func (r *AdminRepository) UpdateLastLogin(id string) error {
	_, err := r.adminDB.Exec("UPDATE admin_users SET last_login_at = $1 WHERE id = $2", time.Now(), id)
	return err
}

func (r *AdminRepository) DeleteAdmin(id string) error {
	_, err := r.adminDB.Exec("DELETE FROM admin_users WHERE id = $1", id)
	return err
}

// ========== Roles ==========

func (r *AdminRepository) GetRoles() ([]models.AdminRole, error) {
	query := `SELECT id, name, description, is_system, created_at FROM admin_roles ORDER BY name`
	rows, err := r.adminDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.AdminRole
	for rows.Next() {
		var role models.AdminRole
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *AdminRepository) GetRoleByID(id string) (*models.AdminRole, error) {
	role := &models.AdminRole{}
	query := `SELECT id, name, description, is_system, created_at FROM admin_roles WHERE id = $1`
	err := r.adminDB.QueryRow(query, id).Scan(&role.ID, &role.Name, &role.Description, &role.IsSystem, &role.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Get permissions
	permQuery := `
		SELECT p.id, p.code, p.name, p.description, p.category
		FROM admin_permissions p
		JOIN admin_role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`
	rows, err := r.adminDB.Query(permQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var perm models.AdminPermission
		if err := rows.Scan(&perm.ID, &perm.Code, &perm.Name, &perm.Description, &perm.Category); err != nil {
			return nil, err
		}
		role.Permissions = append(role.Permissions, perm)
	}

	return role, nil
}

func (r *AdminRepository) GetAdminPermissions(adminID string) ([]string, error) {
	query := `
		SELECT p.code 
		FROM admin_permissions p
		JOIN admin_role_permissions rp ON p.id = rp.permission_id
		JOIN admin_users u ON u.role_id = rp.role_id
		WHERE u.id = $1
	`
	rows, err := r.adminDB.Query(query, adminID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		permissions = append(permissions, code)
	}
	return permissions, nil
}

// ========== Audit Logs ==========

func (r *AdminRepository) CreateAuditLog(log *models.AdminAuditLog) error {
	query := `
		INSERT INTO admin_audit_logs (admin_id, admin_email, action, resource, resource_id, details, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.adminDB.Exec(query,
		log.AdminID, log.AdminEmail, log.Action, log.Resource,
		log.ResourceID, log.Details, log.IPAddress, log.UserAgent,
	)
	return err
}

func (r *AdminRepository) GetAuditLogs(limit, offset int, filters map[string]string) ([]models.AdminAuditLog, error) {
	query := `
		SELECT id, admin_id, admin_email, action, resource, resource_id, details, ip_address, user_agent, created_at
		FROM admin_audit_logs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.adminDB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.AdminAuditLog
	for rows.Next() {
		var log models.AdminAuditLog
		if err := rows.Scan(
			&log.ID, &log.AdminID, &log.AdminEmail, &log.Action,
			&log.Resource, &log.ResourceID, &log.Details,
			&log.IPAddress, &log.UserAgent, &log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

// ========== Main DB Queries (Read-Only) ==========

func (r *AdminRepository) GetUsersFromMainDB(limit, offset int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, is_active, kyc_level, kyc_status, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	return r.queryToMaps(r.mainDB, query, limit, offset)
}

func (r *AdminRepository) GetTransactionsFromMainDB(limit, offset int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, from_wallet_id as wallet_id, transaction_type as type, amount, currency, status, created_at
		FROM transactions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	return r.queryToMaps(r.mainDB, query, limit, offset)
}

func (r *AdminRepository) GetCardsFromMainDB(limit, offset int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, user_id, card_type, status, currency, balance, created_at
		FROM cards
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	return r.queryToMaps(r.mainDB, query, limit, offset)
}

func (r *AdminRepository) GetWalletsFromMainDB(limit, offset int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, user_id, currency, balance, status, created_at
		FROM wallets
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	return r.queryToMaps(r.mainDB, query, limit, offset)
}

func (r *AdminRepository) GetDashboardStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Users count
	var usersCount int
	r.mainDB.QueryRow("SELECT COUNT(*) FROM users").Scan(&usersCount)
	stats["total_users"] = usersCount

	// Active users today
	var activeToday int
	r.mainDB.QueryRow("SELECT COUNT(*) FROM users WHERE last_login_at > CURRENT_DATE").Scan(&activeToday)
	stats["active_today"] = activeToday

	// Transactions today
	var txToday int
	r.mainDB.QueryRow("SELECT COUNT(*) FROM transactions WHERE created_at > CURRENT_DATE").Scan(&txToday)
	stats["transactions_today"] = txToday

	// Total volume today
	var volumeToday float64
	r.mainDB.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE created_at > CURRENT_DATE").Scan(&volumeToday)
	stats["volume_today"] = volumeToday

	// Cards
	var cardsCount int
	r.mainDB.QueryRow("SELECT COUNT(*) FROM cards").Scan(&cardsCount)
	stats["total_cards"] = cardsCount

	// Active cards
	var activeCards int
	r.mainDB.QueryRow("SELECT COUNT(*) FROM cards WHERE status = 'active'").Scan(&activeCards)
	stats["active_cards"] = activeCards

	// Wallets
	var walletsCount int
	r.mainDB.QueryRow("SELECT COUNT(*) FROM wallets").Scan(&walletsCount)
	stats["total_wallets"] = walletsCount

	// Pending KYC - using kyc_status field
	var pendingKYC int
	r.mainDB.QueryRow("SELECT COUNT(*) FROM users WHERE kyc_status = 'pending'").Scan(&pendingKYC)
	stats["pending_kyc"] = pendingKYC

	// Verified KYC
	var verifiedKYC int
	r.mainDB.QueryRow("SELECT COUNT(*) FROM users WHERE kyc_status = 'verified'").Scan(&verifiedKYC)
	stats["verified_kyc"] = verifiedKYC

	// New users today
	var newUsersToday int
	r.mainDB.QueryRow("SELECT COUNT(*) FROM users WHERE created_at > CURRENT_DATE").Scan(&newUsersToday)
	stats["new_users_today"] = newUsersToday

	// Total transfers today
	var transfersToday int
	r.mainDB.QueryRow("SELECT COUNT(*) FROM transfers WHERE created_at > CURRENT_DATE").Scan(&transfersToday)
	stats["transfers_today"] = transfersToday

	return stats, nil
}

func (r *AdminRepository) queryToMaps(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// Convert []byte to string to prevent base64 encoding in JSON
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

// ========== KYC Management (Write to Main DB) ==========

// UpdateUserKYCStatus updates the KYC status of a user in the main database
func (r *AdminRepository) UpdateUserKYCStatus(userID, status string, level int) error {
	query := `UPDATE users SET kyc_status = $1, kyc_level = $2, updated_at = $3 WHERE id = $4`
	_, err := r.mainDB.Exec(query, status, level, time.Now(), userID)
	return err
}

// BlockUser blocks a user in the main database
func (r *AdminRepository) BlockUser(userID string) error {
	query := `UPDATE users SET is_active = false, updated_at = $1 WHERE id = $2`
	_, err := r.mainDB.Exec(query, time.Now(), userID)
	return err
}

// UnblockUser unblocks a user in the main database
func (r *AdminRepository) UnblockUser(userID string) error {
	query := `UPDATE users SET is_active = true, updated_at = $1 WHERE id = $2`
	_, err := r.mainDB.Exec(query, time.Now(), userID)
	return err
}

// UpdateWalletStatus updates a wallet status in the main database
func (r *AdminRepository) UpdateWalletStatus(walletID, status string) error {
	query := `UPDATE wallets SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.mainDB.Exec(query, status, time.Now(), walletID)
	return err
}

// UpdateCardStatus updates a card status in the main database
func (r *AdminRepository) UpdateCardStatus(cardID, status string) error {
	query := `UPDATE cards SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.mainDB.Exec(query, status, time.Now(), cardID)
	return err
}

// GetUserByID gets a user from the main database
func (r *AdminRepository) GetUserByID(userID string) (map[string]interface{}, error) {
	query := `
		SELECT id, email, first_name, last_name, phone, is_active, kyc_level, kyc_status, created_at
		FROM users WHERE id = $1
	`
	rows, err := r.mainDB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := r.rowsToMaps(rows)
	if err != nil || len(results) == 0 {
		return nil, err
	}
	return results[0], nil
}

func (r *AdminRepository) rowsToMaps(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, _ := rows.Columns()
	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// Convert []byte to string to prevent base64 encoding in JSON
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

// GetUserKYCDocuments retrieves all KYC documents for a specific user
func (r *AdminRepository) GetUserKYCDocuments(userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT id, user_id, document_type as type, file_name, file_path, file_size, mime_type, status, 
			file_url, rejection_reason, reviewed_at, reviewed_by, uploaded_at, created_at
		FROM kyc_documents 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`
	return r.queryToMaps(r.mainDB, query, userID)
}

// GetAllKYCRequests retrieves all KYC documents with user info from the main database
func (r *AdminRepository) GetAllKYCRequests(status string, limit, offset int) ([]map[string]interface{}, error) {
	query := `
		SELECT k.id, k.user_id, k.document_type as type, k.file_name, k.file_path, k.file_size, k.mime_type, 
			k.status, k.file_url, k.rejection_reason, k.reviewed_at, k.reviewed_by, k.uploaded_at, k.created_at,
			u.first_name, u.last_name, u.email, u.phone
		FROM kyc_documents k
		JOIN users u ON k.user_id = u.id
	`

	var args []interface{}
	argCount := 1

	if status != "" && status != "all" {
		query += fmt.Sprintf(" WHERE k.status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY k.created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.mainDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse with custom mapping since we joined tables
	columns, _ := rows.Columns()
	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		user := make(map[string]interface{})
		doc := make(map[string]interface{})

		for i, col := range columns {
			val := values[i]
			var v interface{}
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}

			// Map to nested structure if possible, or flat
			if col == "first_name" || col == "last_name" || col == "email" || col == "phone" {
				user[col] = v
			} else {
				doc[col] = v
			}
			row[col] = v // Keep flat for now as existing map util does, or restructure?
			// The frontend expects { document: {...}, user_name: "...", user_email: "..." }
			// Let's rely on the service to restructure, or do it here.
			// Existing "queryToMaps" returns flat map. Let's return flat map and formatting in service?
			// Or just return flat map and let frontend handle it.
			// My Vue code expects: req.user_name, req.user_email, req.document.type
		}

		// Let's manually construct the DTO-like map to match Vue expectation
		// Vue: req.user_name, req.user_email, req.document (object)

		formattedRow := make(map[string]interface{})
		formattedDoc := make(map[string]interface{})

		// Populate doc
		formattedDoc["id"] = row["id"]
		formattedDoc["user_id"] = row["user_id"]
		formattedDoc["type"] = row["type"]
		formattedDoc["file_name"] = row["file_name"]
		formattedDoc["file_path"] = row["file_path"]
		formattedDoc["file_url"] = row["file_url"]
		formattedDoc["status"] = row["status"]
		formattedDoc["uploaded_at"] = row["uploaded_at"]
		formattedDoc["mime_type"] = row["mime_type"]
		// ... add others if needed

		formattedRow["document"] = formattedDoc

		// Populate user info
		fName, _ := row["first_name"].(string)
		lName, _ := row["last_name"].(string)
		formattedRow["user_name"] = fName + " " + lName
		formattedRow["user_email"] = row["email"]
		formattedRow["user_phone"] = row["phone"]

		results = append(results, formattedRow)
	}

	return results, nil
}

// ========== Fee Configuration ==========

func (r *AdminRepository) EnsureFeeConfigsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS fee_configurations (
			id UUID PRIMARY KEY,
			key VARCHAR(100) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			type VARCHAR(50) NOT NULL, -- flat, percentage, hybrid
			fixed_amount DECIMAL(20, 8) DEFAULT 0,
			percentage_amount DECIMAL(10, 4) DEFAULT 0,
			currency VARCHAR(10) DEFAULT 'EUR',
			is_enabled BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			updated_by VARCHAR(255)
		);
	`
	_, err := r.adminDB.Exec(query)
	return err
}

func (r *AdminRepository) GetFeeConfigs() ([]models.FeeConfig, error) {
	// Query to get all configs
	query := `SELECT id, key, name, description, type, fixed_amount, percentage_amount, currency, is_enabled, created_at, updated_at, updated_by FROM fee_configurations ORDER BY key`

	rows, err := r.adminDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []models.FeeConfig
	for rows.Next() {
		var c models.FeeConfig
		var desc, upBy sql.NullString

		err := rows.Scan(
			&c.ID, &c.Key, &c.Name, &desc, &c.Type,
			&c.FixedAmount, &c.PercentageAmount, &c.Currency, &c.IsEnabled,
			&c.CreatedAt, &c.UpdatedAt, &upBy,
		)
		if err != nil {
			return nil, err
		}

		if desc.Valid {
			c.Description = desc.String
		}
		if upBy.Valid {
			c.UpdatedBy = upBy.String
		}

		configs = append(configs, c)
	}
	return configs, nil
}

func (r *AdminRepository) GetFeeConfigByKey(key string) (*models.FeeConfig, error) {
	query := `SELECT id, key, name, description, type, fixed_amount, percentage_amount, currency, is_enabled, created_at, updated_at, updated_by FROM fee_configurations WHERE key = $1`

	var c models.FeeConfig
	var desc, upBy sql.NullString

	err := r.adminDB.QueryRow(query, key).Scan(
		&c.ID, &c.Key, &c.Name, &desc, &c.Type,
		&c.FixedAmount, &c.PercentageAmount, &c.Currency, &c.IsEnabled,
		&c.CreatedAt, &c.UpdatedAt, &upBy,
	)
	if err != nil {
		return nil, err
	}

	if desc.Valid {
		c.Description = desc.String
	}
	if upBy.Valid {
		c.UpdatedBy = upBy.String
	}

	return &c, nil
}

func (r *AdminRepository) UpdateFeeConfig(key string, req models.UpdateFeeRequest, updatedBy string) error {
	query := `UPDATE fee_configurations SET updated_at = $1, updated_by = $2`
	args := []interface{}{time.Now(), updatedBy}
	argCount := 3

	if req.Type != nil {
		query += fmt.Sprintf(", type = $%d", argCount)
		args = append(args, *req.Type)
		argCount++
	}
	if req.FixedAmount != nil {
		query += fmt.Sprintf(", fixed_amount = $%d", argCount)
		args = append(args, *req.FixedAmount)
		argCount++
	}
	if req.PercentageAmount != nil {
		query += fmt.Sprintf(", percentage_amount = $%d", argCount)
		args = append(args, *req.PercentageAmount)
		argCount++
	}
	if req.IsEnabled != nil {
		query += fmt.Sprintf(", is_enabled = $%d", argCount)
		args = append(args, *req.IsEnabled)
		argCount++
	}

	query += fmt.Sprintf(" WHERE key = $%d", argCount)
	args = append(args, key)

	_, err := r.adminDB.Exec(query, args...)
	return err
}

func (r *AdminRepository) CreateFeeConfig(c *models.FeeConfig) error {
	query := `
		INSERT INTO fee_configurations (id, key, name, description, type, fixed_amount, percentage_amount, currency, is_enabled, created_at, updated_at, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.adminDB.Exec(query,
		c.ID, c.Key, c.Name, c.Description, c.Type,
		c.FixedAmount, c.PercentageAmount, c.Currency, c.IsEnabled,
		c.CreatedAt, c.UpdatedAt, c.UpdatedBy,
	)
	return err
}
