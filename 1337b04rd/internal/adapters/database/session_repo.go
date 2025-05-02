package database

import (
	"1337b04rd/internal/app/domain/models"
	"database/sql"
)

type PostgresSessionRepo struct {
	DB *sql.DB
}

func (r *PostgresSessionRepo) GetSessionData(sessionID string) (models.UserData, bool) {
	var data models.UserData
	row := r.DB.QueryRow("SELECT name, avatar, last_visit FROM sessions WHERE id=$1", sessionID)
	err := row.Scan(&data.Name, &data.Avatar, &data.LastVisit)
	if err != nil {
		return models.UserData{}, false
	}
	return data, true
}

func (r *PostgresSessionRepo) SetSessionData(sessionID string, data models.UserData) error {
	_, err := r.DB.Exec(
		`INSERT INTO sessions (id, name, avatar) 
		 VALUES ($1, $2, $3)
		 ON CONFLICT (id) DO UPDATE SET name = $2, avatar = $3`,
		sessionID, data.Name, data.Avatar,
	)
	return err
}
