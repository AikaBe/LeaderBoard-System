package api

import (
	"1337b04rd/internal/app/domain/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

var UserVisitData = make(map[string]models.UserData)

type Character struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

func GetRandomCharacter() (Character, error) {
	id := rand.Intn(826) + 1
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
