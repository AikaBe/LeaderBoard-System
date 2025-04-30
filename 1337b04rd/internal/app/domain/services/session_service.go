package services

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

type SessionService struct {
	Repo ports.SessionRepository
}

func NewSessionService(repo ports.SessionRepository) *SessionService {
	return &SessionService{Repo: repo}
}

func (s *SessionService) GenerateSessionID() string {
	rand.Seed(time.Now().UnixNano())
	var id [16]byte
	for i := range id {
		id[i] = byte(rand.Intn(256))
	}
	return hex.EncodeToString(id[:])
}

func (s *SessionService) CreateOrUpdateSession(sessionID string, name, avatar string) string {
	now := time.Now().Format(time.RFC3339)
	if sessionID == "" {
		sessionID = s.GenerateSessionID()
		fmt.Println("Генерация нового sessionId:", sessionID)
	}
	s.Repo.SetSessionData(sessionID, models.UserData{
		Name:      name,
		Avatar:    avatar,
		LastVisit: now,
	})
	fmt.Println("Сессия сохранена для:", name)
	return sessionID
}

func (s *SessionService) GetUserData(sessionID string) (models.UserData, bool) {
	return s.Repo.GetSessionData(sessionID)
}
