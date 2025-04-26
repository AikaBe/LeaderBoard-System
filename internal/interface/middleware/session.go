package middleware

import (
	"1337b04rd/internal/adapters/api"
	"1337b04rd/internal/app/domain/services"
	"net/http"
	"time"
)

var sessionService = services.NewSessionService(&api.InMemorySessionRepo{})

func getCookieValue(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func LoginOrLastVisitHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId, _ := getCookieValue(r, "sessionId")

		if sessionId == "" {
			character, err := api.GetNextCharacter()
			if err != nil {
				http.Error(w, "error with character!", http.StatusInternalServerError)
				return
			}
			sessionId = sessionService.CreateOrUpdateSession("", character.Name, character.Image)
			http.SetCookie(w, &http.Cookie{
				Name:     "sessionId",
				Value:    sessionId,
				Path:     "/",
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				HttpOnly: true,
			})
		} else {
			userData, ok := sessionService.GetUserData(sessionId)
			if !ok {
				http.Error(w, "Session not found", http.StatusNotFound)
				return
			}

			userData.LastVisit = time.Now().Format(time.RFC3339)
			sessionService.CreateOrUpdateSession(sessionId, userData.Name, userData.Avatar)
		}
		next.ServeHTTP(w, r)
	})
}
