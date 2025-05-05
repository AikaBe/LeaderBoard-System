package ports

import "1337b04rd/internal/app/domain/models"

type SessionRepository interface {
	GetSessionData(sessionID string) (models.UserData, bool)
	SetSessionData(sessionID string, data models.UserData) error
}
