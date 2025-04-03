package main

import (
	"1337b04rd/internal/interface/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/setcookie", handlers.SetCookie)
	http.HandleFunc("/getcookie", handlers.GetCookie)

	fmt.Println("Server start at http://localhost:8080 ")
	http.ListenAndServe(":8080", nil)
}
