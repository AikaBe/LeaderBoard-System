package api

import (
	"1337b04rd/internal/app/domain/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var (
	UserVisitData = make(map[string]models.UserData)
	currentIndex  = 1
	maxIndex      = 826
	mu            sync.Mutex
)

type Character struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func GetNextCharacter() (Character, error) {
	mu.Lock()
	id := currentIndex
	currentIndex++
	if currentIndex > maxIndex {
		currentIndex = 1
	}
	mu.Unlock()

	url := fmt.Sprintf("https://rickandmortyapi.com/api/character/%d", id)

	resp, err := http.Get(url)
	if err != nil {
		return Character{}, err
	}
	defer resp.Body.Close()

	var character Character
	err = json.NewDecoder(resp.Body).Decode(&character)
	if err != nil {
		return Character{}, err
	}

	return character, nil
}
