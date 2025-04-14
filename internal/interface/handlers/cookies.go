package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, r *http.Request) {
	// create cookie
	cookie := &http.Cookie{
		Name:     "ID",
		Value:    "John",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
	}

	// set the cookie as an answer
	http.SetCookie(w, cookie)

	fmt.Fprintf(w, "Cookie 'Expires' installed: %s", cookie.Expires)
}

func GetCookie(w http.ResponseWriter, r *http.Request) {
	// read cookie
	cookie, err := r.Cookie("ID")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}

	// return the answer
	fmt.Fprintf(w, "Cookie 'ID' have the value : %s", cookie.Value)
}
