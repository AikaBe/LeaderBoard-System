package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"1337b04rd/internal/app/domain/models"
)

var (
	UserVisitData = make(map[string]models.UserData)
	currentIndex  = 1
	maxIndex      = 826
	mu            sync.Mutex
)

// Character represents the structure of a character from the API.
type Character struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// GetNextCharacter retrieves the next character from the API.
func GetNextCharacter() (Character, error) {
	mu.Lock()
	id := currentIndex
	currentIndex++
	if currentIndex > maxIndex {
		currentIndex = 1
	}
	mu.Unlock()

	// Format the URL for the API request
	url := fmt.Sprintf("https://rickandmortyapi.com/api/character/%d", id)

	// Log the request
	slog.Info("Fetching character", "id", id)

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Failed to fetch character", "error", err)
		return Character{}, err
	}
	defer resp.Body.Close()

	var character Character
	err = json.NewDecoder(resp.Body).Decode(&character)
	if err != nil {
		slog.Error("Failed to decode character response", "error", err)
		return Character{}, err
	}

	// Log successful response
	slog.Info("Successfully fetched character", "name", character.Name)

	return character, nil
}
