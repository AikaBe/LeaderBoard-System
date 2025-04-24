package middleware

import (
	"1337b04rd/internal/adapters/api"
	"1337b04rd/internal/app/domain/services"
	"fmt"
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

func LoginOrLastVisitHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем sessionId из cookies
	sessionId, _ := getCookieValue(r, "sessionId")

	if sessionId == "" {
		// Если сессии нет, создаем новую сессию
		character, err := api.GetNextCharacter()
		if err != nil {
			http.Error(w, "Ошибка при получении персонажа", http.StatusInternalServerError)
			return
		}

		// Создаем или обновляем сессию
		sessionId = sessionService.CreateOrUpdateSession("", character.Name, character.Image)

		// Устанавливаем cookie с sessionId
		http.SetCookie(w, &http.Cookie{
			Name:     "sessionId",
			Value:    sessionId,
			Path:     "/",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HttpOnly: true,
		})

		// Отправляем пользователю приветствие и аватар
		fmt.Fprintf(w, "Создана новая сессия для %s<br><img src='%s'>", character.Name, character.Image)
		return
	}

	// Если сессия существует, выводим информацию о последнем визите
	userData, ok := sessionService.GetUserData(sessionId)
	if !ok {
		http.Error(w, "Сессия не найдена", http.StatusUnauthorized)
		return
	}

	// Обновляем дату последнего визита
	userData.LastVisit = time.Now().Format(time.RFC3339)
	sessionService.CreateOrUpdateSession(sessionId, userData.Name, userData.Avatar)

	// Отправляем приветствие и информацию о последнем визите
	fmt.Fprintf(w, "Добро пожаловать назад, %s! Последний визит: %s<br><img src='%s'>", userData.Name, userData.LastVisit, userData.Avatar)
}
