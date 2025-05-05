package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func SetCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "ID",
		Value:    "John",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	slog.Info("Cookie set", "name", cookie.Name, "value", cookie.Value, "expires", cookie.Expires)

	fmt.Fprintf(w, "Cookie 'ID' set with expiration: %s", cookie.Expires)
}

func GetCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("ID")
	if err != nil {
		slog.Warn("Cookie not found", "name", "ID", "error", err)
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}

	slog.Info("Cookie retrieved", "name", cookie.Name, "value", cookie.Value)

	fmt.Fprintf(w, "Cookie 'ID' has value: %s", cookie.Value)
}
