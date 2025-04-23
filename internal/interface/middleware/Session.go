package middleware

/*Проверить, есть ли у пользователя куки с идентификатором сессии.

Если куки есть, использовать этот идентификатор для извлечения данных о последнем визите из базы данных или другого хранилища.

Если куки нет, создать новую сессию, назначить куки и сохранить данные о визите.*/

import (
	"1337b04rd/internal/adapters/api"
	"1337b04rd/internal/app/domain/models"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Эмуляция базы данных для хранения данных о визите

// Генерация уникального идентификатора для сессии
func GenerateSessionID() string {
	// Генерация случайного ID для сессии
	rand.Seed(time.Now().UnixNano())
	var id [16]byte
	for i := range id {
		id[i] = byte(rand.Intn(256))
	}
	return hex.EncodeToString(id[:])
}

// Проверка наличия куки в запросе
func GetCookieValue(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Обработчик входа пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, err := GetCookieValue(r, "sessionId")
	if err != nil || sessionId == "" {

		sessionId = GenerateSessionID()
		http.SetCookie(w, &http.Cookie{
			Name:     "sessionId",
			Value:    sessionId,
			Path:     "/",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			HttpOnly: true,
		})
		character, err := api.GetNextCharacter()
		if err != nil {
			http.Error(w, "Ошибка при получении персонажа", http.StatusInternalServerError)
			return
		}

		api.UserVisitData[sessionId] = models.UserData{
			LastVisit: time.Now().Format(time.RFC3339),
			Name:      character.Name,
			Avatar:    character.Image,
		}

		fmt.Fprintf(w, "Создана новая сессия для %s<br><img src='%s'>", character.Name, character.Image)
	} else {
		data := api.UserVisitData[sessionId]
		data.LastVisit = time.Now().Format(time.RFC3339)
		api.UserVisitData[sessionId] = data

		fmt.Fprintf(w, "Добро пожаловать назад, %s!<br><img src='%s'>", data.Name, data.Avatar)
	}
}

// Обработчик для получения последнего визита пользователя
func LastVisitHandler(w http.ResponseWriter, r *http.Request) {
	sessionId, err := GetCookieValue(r, "sessionId")
	if err != nil || sessionId == "" {
		http.Error(w, "Не найдена сессия", http.StatusUnauthorized)
		return
	}

	data, exists := api.UserVisitData[sessionId]
	if !exists {
		http.Error(w, "Нет данных о визите", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Привет, %s! Последний визит: %s<br><img src='%s'>", data.Name, data.LastVisit, data.Avatar)
}
