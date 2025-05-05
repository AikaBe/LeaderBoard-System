package api

import (
	"database/sql"
	"errors"
	"log/slog"

	"1337b04rd/internal/app/domain/models"
)

// DBSessionRepo represents a repository for handling session data.
type DBSessionRepo struct {
	DB *sql.DB
}

// NewDBSessionRepo creates a new instance of DBSessionRepo.
func NewDBSessionRepo(db *sql.DB) *DBSessionRepo {
	return &DBSessionRepo{DB: db}
}

// GetSessionData retrieves the session data for the given session ID.
func (r *DBSessionRepo) GetSessionData(sessionID string) (models.UserData, bool) {
	var data models.UserData
	query := `SELECT name, avatar FROM sessions WHERE session_id = $1`
	row := r.DB.QueryRow(query, sessionID)

	err := row.Scan(&data.Name, &data.Avatar)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No session found, log this case
			slog.Info("No session found", "sessionID", sessionID)
			return models.UserData{}, false
		}
		// Log error retrieving session data
		slog.Error("Error retrieving session data", "sessionID", sessionID, "error", err)
		return models.UserData{}, false
	}

	// Successfully retrieved session data
	slog.Info("Successfully retrieved session data", "sessionID", sessionID, "name", data.Name)
	return data, true
}

// SetSessionData stores session data for a given session ID.
func (r *DBSessionRepo) SetSessionData(sessionID string, data models.UserData) error {
	query := `
		INSERT INTO sessions (session_id, name, avatar)
		VALUES ($1, $2, $3)
		ON CONFLICT (session_id) DO UPDATE
		SET name = EXCLUDED.name, avatar = EXCLUDED.avatar;
	`
	_, err := r.DB.Exec(query, sessionID, data.Name, data.Avatar)
	if err != nil {
		// Log error saving session data
		slog.Error("Error saving session data", "sessionID", sessionID, "error", err)
	}
	// Log successful save operation
	slog.Info("Session data saved or updated", "sessionID", sessionID)
	return err
}
