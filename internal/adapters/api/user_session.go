package api

import "1337b04rd/internal/app/domain/models"

type InMemorySessionRepo struct{}

func (r *InMemorySessionRepo) GetSessionData(sessionID string) (models.UserData, bool) {
	mu.Lock()
	defer mu.Unlock()
	data, ok := UserVisitData[sessionID]
	return data, ok
}

func (r *InMemorySessionRepo) SetSessionData(sessionID string, data models.UserData) {
	mu.Lock()
	defer mu.Unlock()
	UserVisitData[sessionID] = data
}
