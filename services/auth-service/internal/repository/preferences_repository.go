package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/crypto-bank/microservices-financial-app/services/auth-service/internal/models"
)

type PreferencesRepository struct {
	db *sql.DB
}

func NewPreferencesRepository(db *sql.DB) *PreferencesRepository {
	return &PreferencesRepository{db: db}
}

// ============ User Preferences ============

func (r *PreferencesRepository) GetUserPreferences(userID string) (*models.UserPreferences, error) {
	query := `SELECT id, user_id, theme, language, currency, timezone, number_format, date_format, 
		hide_balances, analytics_enabled, auto_lock_minutes, created_at, updated_at 
		FROM user_preferences WHERE user_id = $1`

	var prefs models.UserPreferences
	err := r.db.QueryRow(query, userID).Scan(
		&prefs.ID, &prefs.UserID, &prefs.Theme, &prefs.Language, &prefs.Currency, &prefs.Timezone,
		&prefs.NumberFormat, &prefs.DateFormat, &prefs.HideBalances, &prefs.AnalyticsEnabled,
		&prefs.AutoLockMinutes, &prefs.CreatedAt, &prefs.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return r.CreateDefaultPreferences(userID)
	}
	if err != nil {
		return nil, err
	}
	return &prefs, nil
}

func (r *PreferencesRepository) CreateDefaultPreferences(userID string) (*models.UserPreferences, error) {
	return r.CreateDefaultPreferencesWithCountry(userID, "")
}

// CreateDefaultPreferencesWithCountry creates default preferences with currency based on country
func (r *PreferencesRepository) CreateDefaultPreferencesWithCountry(userID, countryCode string) (*models.UserPreferences, error) {
	currency := getCurrencyForCountry(countryCode)
	
	prefs := &models.UserPreferences{
		ID:               uuid.New().String(),
		UserID:           userID,
		Theme:            "dark",
		Language:         getLanguageForCountry(countryCode),
		Currency:         currency,
		Timezone:         getTimezoneForCountry(countryCode),
		NumberFormat:     "fr",
		DateFormat:       "DD/MM/YYYY",
		HideBalances:     false,
		AnalyticsEnabled: true,
		AutoLockMinutes:  5,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err := r.db.Exec(`
		INSERT INTO user_preferences (id, user_id, theme, language, currency, timezone, 
			number_format, date_format, hide_balances, analytics_enabled, auto_lock_minutes, 
			created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, prefs.ID, prefs.UserID, prefs.Theme, prefs.Language, prefs.Currency, prefs.Timezone,
		prefs.NumberFormat, prefs.DateFormat, prefs.HideBalances, prefs.AnalyticsEnabled,
		prefs.AutoLockMinutes, prefs.CreatedAt, prefs.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return prefs, nil
}

// getCurrencyForCountry returns the appropriate currency for a country code
func getCurrencyForCountry(countryCode string) string {
	// UEMOA - West African CFA Franc (XOF)
	xofCountries := map[string]bool{
		"BJ": true, // Benin
		"BF": true, // Burkina Faso
		"CI": true, // Côte d'Ivoire
		"GW": true, // Guinea-Bissau
		"ML": true, // Mali
		"NE": true, // Niger
		"SN": true, // Senegal
		"TG": true, // Togo
	}

	// CEMAC - Central African CFA Franc (XAF)
	xafCountries := map[string]bool{
		"CM": true, // Cameroon
		"CF": true, // Central African Republic
		"TD": true, // Chad
		"CG": true, // Congo
		"GQ": true, // Equatorial Guinea
		"GA": true, // Gabon
	}

	// Eurozone
	eurCountries := map[string]bool{
		"FR": true, "DE": true, "IT": true, "ES": true, "PT": true,
		"NL": true, "BE": true, "AT": true, "IE": true, "FI": true,
		"GR": true, "SK": true, "SI": true, "EE": true, "LV": true,
		"LT": true, "CY": true, "MT": true, "LU": true,
	}

	// UK
	if countryCode == "GB" || countryCode == "UK" {
		return "GBP"
	}

	// USA and territories
	if countryCode == "US" || countryCode == "PR" || countryCode == "VI" {
		return "USD"
	}

	// Canada
	if countryCode == "CA" {
		return "CAD"
	}

	// Check currency zones
	if xofCountries[countryCode] {
		return "XOF"
	}
	if xafCountries[countryCode] {
		return "XAF"
	}
	if eurCountries[countryCode] {
		return "EUR"
	}

	// Other African countries
	switch countryCode {
	case "MA":
		return "MAD" // Morocco
	case "DZ":
		return "DZD" // Algeria
	case "TN":
		return "TND" // Tunisia
	case "EG":
		return "EGP" // Egypt
	case "NG":
		return "NGN" // Nigeria
	case "GH":
		return "GHS" // Ghana
	case "KE":
		return "KES" // Kenya
	case "ZA":
		return "ZAR" // South Africa
	case "TZ":
		return "TZS" // Tanzania
	case "UG":
		return "UGX" // Uganda
	case "RW":
		return "RWF" // Rwanda
	case "CD":
		return "CDF" // DR Congo
	case "AO":
		return "AOA" // Angola
	case "MZ":
		return "MZN" // Mozambique
	case "ET":
		return "ETB" // Ethiopia
	case "GN":
		return "GNF" // Guinea
	}

	// Default to USD for unknown countries
	return "USD"
}

// getLanguageForCountry returns appropriate language
func getLanguageForCountry(countryCode string) string {
	frenchCountries := map[string]bool{
		"FR": true, "BE": true, "CH": true, "CA": true,
		"BJ": true, "BF": true, "CI": true, "GW": true, "ML": true,
		"NE": true, "SN": true, "TG": true, "CM": true, "CF": true,
		"TD": true, "CG": true, "GA": true, "MA": true, "DZ": true,
		"TN": true, "CD": true, "RW": true, "MG": true,
	}
	if frenchCountries[countryCode] {
		return "fr"
	}
	return "en"
}

// getTimezoneForCountry returns appropriate timezone
func getTimezoneForCountry(countryCode string) string {
	switch countryCode {
	case "US":
		return "America/New_York"
	case "GB", "UK":
		return "Europe/London"
	case "FR", "DE", "IT", "ES", "BE", "NL":
		return "Europe/Paris"
	case "CI", "SN", "ML", "BF", "NE", "TG", "BJ", "GH", "GN":
		return "Africa/Abidjan"
	case "CM", "GA", "CG", "CF", "TD", "NG":
		return "Africa/Lagos"
	case "MA":
		return "Africa/Casablanca"
	case "EG":
		return "Africa/Cairo"
	case "ZA":
		return "Africa/Johannesburg"
	case "KE", "TZ", "UG", "RW", "ET":
		return "Africa/Nairobi"
	default:
		return "UTC"
	}
}

func (r *PreferencesRepository) UpdateUserPreferences(userID string, req *models.UpdatePreferencesRequest) (*models.UserPreferences, error) {
	prefs, err := r.GetUserPreferences(userID)
	if err != nil {
		return nil, err
	}

	if req.Theme != nil {
		prefs.Theme = *req.Theme
	}
	if req.Language != nil {
		prefs.Language = *req.Language
	}
	if req.Currency != nil {
		prefs.Currency = *req.Currency
	}
	if req.Timezone != nil {
		prefs.Timezone = *req.Timezone
	}
	if req.NumberFormat != nil {
		prefs.NumberFormat = *req.NumberFormat
	}
	if req.DateFormat != nil {
		prefs.DateFormat = *req.DateFormat
	}
	if req.HideBalances != nil {
		prefs.HideBalances = *req.HideBalances
	}
	if req.AnalyticsEnabled != nil {
		prefs.AnalyticsEnabled = *req.AnalyticsEnabled
	}
	if req.AutoLockMinutes != nil {
		prefs.AutoLockMinutes = *req.AutoLockMinutes
	}

	prefs.UpdatedAt = time.Now()

	_, err = r.db.Exec(`
		UPDATE user_preferences SET theme = $1, language = $2, currency = $3, timezone = $4,
			number_format = $5, date_format = $6, hide_balances = $7, analytics_enabled = $8,
			auto_lock_minutes = $9, updated_at = $10
		WHERE user_id = $11
	`, prefs.Theme, prefs.Language, prefs.Currency, prefs.Timezone,
		prefs.NumberFormat, prefs.DateFormat, prefs.HideBalances, prefs.AnalyticsEnabled,
		prefs.AutoLockMinutes, prefs.UpdatedAt, userID)

	return prefs, err
}

// ============ Notification Preferences ============

func (r *PreferencesRepository) GetNotificationPrefs(userID string) (*models.NotificationPreferences, error) {
	query := `SELECT id, user_id, push_enabled, transfer_received, transfer_sent, card_payment, 
		low_balance, new_login, password_change, otp_via_sms, weekly_report, newsletter, 
		promotions, created_at, updated_at 
		FROM notification_preferences WHERE user_id = $1`

	var prefs models.NotificationPreferences
	err := r.db.QueryRow(query, userID).Scan(
		&prefs.ID, &prefs.UserID, &prefs.PushEnabled, &prefs.TransferReceived, &prefs.TransferSent,
		&prefs.CardPayment, &prefs.LowBalance, &prefs.NewLogin, &prefs.PasswordChange,
		&prefs.OtpViaSMS, &prefs.WeeklyReport, &prefs.Newsletter, &prefs.Promotions,
		&prefs.CreatedAt, &prefs.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return r.CreateDefaultNotificationPrefs(userID)
	}
	if err != nil {
		return nil, err
	}
	return &prefs, nil
}

func (r *PreferencesRepository) CreateDefaultNotificationPrefs(userID string) (*models.NotificationPreferences, error) {
	prefs := &models.NotificationPreferences{
		ID:               uuid.New().String(),
		UserID:           userID,
		PushEnabled:      true,
		TransferReceived: true,
		TransferSent:     true,
		CardPayment:      true,
		LowBalance:       true,
		NewLogin:         true,
		PasswordChange:   true,
		OtpViaSMS:        true,
		WeeklyReport:     false,
		Newsletter:       false,
		Promotions:       false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	_, err := r.db.Exec(`
		INSERT INTO notification_preferences (id, user_id, push_enabled, transfer_received, 
			transfer_sent, card_payment, low_balance, new_login, password_change, otp_via_sms,
			weekly_report, newsletter, promotions, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`, prefs.ID, prefs.UserID, prefs.PushEnabled, prefs.TransferReceived,
		prefs.TransferSent, prefs.CardPayment, prefs.LowBalance, prefs.NewLogin,
		prefs.PasswordChange, prefs.OtpViaSMS, prefs.WeeklyReport, prefs.Newsletter,
		prefs.Promotions, prefs.CreatedAt, prefs.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return prefs, nil
}

func (r *PreferencesRepository) UpdateNotificationPrefs(userID string, req *models.UpdateNotificationPrefsRequest) (*models.NotificationPreferences, error) {
	prefs, err := r.GetNotificationPrefs(userID)
	if err != nil {
		return nil, err
	}

	if req.PushEnabled != nil {
		prefs.PushEnabled = *req.PushEnabled
	}
	if req.TransferReceived != nil {
		prefs.TransferReceived = *req.TransferReceived
	}
	if req.TransferSent != nil {
		prefs.TransferSent = *req.TransferSent
	}
	if req.CardPayment != nil {
		prefs.CardPayment = *req.CardPayment
	}
	if req.LowBalance != nil {
		prefs.LowBalance = *req.LowBalance
	}
	if req.NewLogin != nil {
		prefs.NewLogin = *req.NewLogin
	}
	if req.PasswordChange != nil {
		prefs.PasswordChange = *req.PasswordChange
	}
	if req.OtpViaSMS != nil {
		prefs.OtpViaSMS = *req.OtpViaSMS
	}
	if req.WeeklyReport != nil {
		prefs.WeeklyReport = *req.WeeklyReport
	}
	if req.Newsletter != nil {
		prefs.Newsletter = *req.Newsletter
	}
	if req.Promotions != nil {
		prefs.Promotions = *req.Promotions
	}

	prefs.UpdatedAt = time.Now()

	_, err = r.db.Exec(`
		UPDATE notification_preferences SET push_enabled = $1, transfer_received = $2,
			transfer_sent = $3, card_payment = $4, low_balance = $5, new_login = $6,
			password_change = $7, otp_via_sms = $8, weekly_report = $9, newsletter = $10,
			promotions = $11, updated_at = $12
		WHERE user_id = $13
	`, prefs.PushEnabled, prefs.TransferReceived, prefs.TransferSent, prefs.CardPayment,
		prefs.LowBalance, prefs.NewLogin, prefs.PasswordChange, prefs.OtpViaSMS,
		prefs.WeeklyReport, prefs.Newsletter, prefs.Promotions, prefs.UpdatedAt, userID)

	return prefs, err
}

// ============ KYC Documents ============

// UpdateUserKYCStatus updates the user's KYC status after document upload
func (r *PreferencesRepository) UpdateUserKYCStatus(userID, status string) error {
	_, err := r.db.Exec(`
		UPDATE users SET kyc_status = $1, updated_at = $2 WHERE id = $3
	`, status, time.Now(), userID)
	return err
}

func (r *PreferencesRepository) CreateKYCDocument(doc *models.KYCDocument, fileURL string) error {
	doc.ID = uuid.New().String()
	doc.Status = "pending"
	doc.UploadedAt = time.Now()
	doc.CreatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO kyc_documents (id, user_id, document_type, file_name, file_path, file_size, 
			mime_type, document_number, expiry_date, status, uploaded_at, created_at, file_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, doc.ID, doc.UserID, doc.Type, doc.FileName, doc.FilePath, doc.FileSize,
		doc.MimeType, doc.DocumentNumber, doc.ExpiryDate, doc.Status, doc.UploadedAt, doc.CreatedAt, fileURL)

	return err
}

func (r *PreferencesRepository) GetKYCDocuments(userID string) ([]models.KYCDocument, error) {
	query := `SELECT id, user_id, document_type, file_name, file_path, file_size, mime_type, 
		document_number, expiry_date, status, rejection_reason, reviewed_at, reviewed_by, uploaded_at, created_at 
		FROM kyc_documents WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []models.KYCDocument
	for rows.Next() {
		var doc models.KYCDocument
		err := rows.Scan(
			&doc.ID, &doc.UserID, &doc.Type, &doc.FileName, &doc.FilePath, &doc.FileSize,
			&doc.MimeType, &doc.DocumentNumber, &doc.ExpiryDate, &doc.Status, &doc.RejectionReason, 
			&doc.ReviewedAt, &doc.ReviewedBy, &doc.UploadedAt, &doc.CreatedAt,
		)
		if err != nil {
			continue
		}
		docs = append(docs, doc)
	}

	if docs == nil {
		docs = []models.KYCDocument{}
	}
	return docs, nil
}

