package services

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
	"encoding/hex"
	"log/slog"
	"math/rand"
	"time"
)

// SessionService manages user sessions.
type SessionService struct {
	Repo ports.SessionRepository
}

// NewSessionService creates a new SessionService.
func NewSessionService(repo ports.SessionRepository) *SessionService {
	return &SessionService{Repo: repo}
}

// GenerateSessionID creates a new random session ID.
func (s *SessionService) GenerateSessionID() string {
	rand.Seed(time.Now().UnixNano())
	var id [16]byte
	for i := range id {
		id[i] = byte(rand.Intn(256))
	}
	return hex.EncodeToString(id[:])
}

// CreateOrUpdateSession sets session data for a user or updates an existing one.
func (s *SessionService) CreateOrUpdateSession(sessionID, name, avatar string) string {
	now := time.Now().Format(time.RFC3339)
	if sessionID == "" {
		sessionID = s.GenerateSessionID()
		slog.Info("Generated new session ID", "sessionID", sessionID)
	}
	s.Repo.SetSessionData(sessionID, models.UserData{
		Name:      name,
		Avatar:    avatar,
		LastVisit: now,
	})
	slog.Info("Session saved", "sessionID", sessionID, "user", name)
	return sessionID
}

// GetUserData retrieves session data for a given session ID.
func (s *SessionService) GetUserData(sessionID string) (models.UserData, bool) {
	userData, ok := s.Repo.GetSessionData(sessionID)
	if !ok {
		slog.Warn("Session not found", "sessionID", sessionID)
	}
	return userData, ok
}
