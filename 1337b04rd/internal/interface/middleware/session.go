package middleware

import (
	"1337b04rd/internal/adapters/api"
	"1337b04rd/internal/app/domain/services"
	"context"
	"log"
	"net/http"
	"time"
)

type AuthMiddleware struct {
	SessionService *services.SessionService
}

// getCookieValue retrieves the value of a cookie by its name
func getCookieValue(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		// Log error when getting the cookie
		log.Printf("error getting cookie '%s': %v", cookieName, err)
		return "", err
	}
	return cookie.Value, nil
}

// LoginOrLastVisitHandler checks the user's session and updates the last visit information
func (am *AuthMiddleware) LoginOrLastVisitHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get the session ID from the cookie
		sessionId, err := getCookieValue(r, "sessionId")
		if sessionId == "" || err != nil {
			// If the session doesn't exist, create a new character
			log.Println("No session found, creating new session with character")
			character, err := api.GetNextCharacter()
			if err != nil {
				// Log error when getting character
				log.Printf("error with character: %v", err)
				http.Error(w, "error with character!", http.StatusInternalServerError)
				return
			}

			// Create a new session
			sessionId = am.SessionService.CreateOrUpdateSession("", character.Name, character.Image)
			log.Printf("New session created for character: %s", character.Name)

			// Set the session ID in the cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "sessionId",
				Value:    sessionId,
				Path:     "/",
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				HttpOnly: true,
			})
			log.Printf("Session ID '%s' set in cookie", sessionId)
		} else {
			// If the session exists, update the last visit data
			log.Printf("Session '%s' found, updating last visit", sessionId)
			userData, ok := am.SessionService.GetUserData(sessionId)
			if !ok {
				log.Printf("session '%s' not found", sessionId)
				http.Error(w, "Session not found", http.StatusNotFound)
				return
			}

			// Update the last visit time
			userData.LastVisit = time.Now().Format(time.RFC3339)
			am.SessionService.CreateOrUpdateSession(sessionId, userData.Name, userData.Avatar)
			log.Printf("Session '%s' updated with new last visit time", sessionId)
		}

		// Add sessionId to the request context
		ctx := context.WithValue(r.Context(), "sessionId", sessionId)
		log.Printf("Passing sessionId '%s' to next handler", sessionId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
