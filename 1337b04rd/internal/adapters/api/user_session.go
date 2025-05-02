package api

import (
	"1337b04rd/internal/app/domain/models"
	"database/sql"
	"errors"
	"log"
)

type DBSessionRepo struct {
	DB *sql.DB
}

func NewDBSessionRepo(db *sql.DB) *DBSessionRepo {
	return &DBSessionRepo{DB: db}
}

func (r *DBSessionRepo) GetSessionData(sessionID string) (models.UserData, bool) {
	var data models.UserData
	query := `SELECT name, avatar FROM sessions WHERE session_id = $1`
	row := r.DB.QueryRow(query, sessionID)

	err := row.Scan(&data.Name, &data.Avatar)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserData{}, false
		}
		log.Println("error retrieving session data:", err)
		return models.UserData{}, false
	}

	return data, true
}

func (r *DBSessionRepo) SetSessionData(sessionID string, data models.UserData) error {
	query := `
		INSERT INTO sessions (session_id, name, avatar)
		VALUES ($1, $2, $3)
		ON CONFLICT (session_id) DO UPDATE
		SET name = EXCLUDED.name, avatar = EXCLUDED.avatar;
	`
	_, err := r.DB.Exec(query, sessionID, data.Name, data.Avatar)
	if err != nil {
		log.Println("error saving session data:", err)
	}
	return err
}
