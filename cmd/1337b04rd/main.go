package main

import (
	"1337b04rd/internal/interface/middleware"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", middleware.LoginHandler)
	http.HandleFunc("/last-visit", middleware.LastVisitHandler)

	// Запуск сервера
	log.Println("Сервер работает на порту 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
	}
}
