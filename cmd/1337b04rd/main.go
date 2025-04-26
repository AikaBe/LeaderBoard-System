package main

import (
	"1337b04rd/internal/interface/middleware"
	"database/sql"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

func main() {
	// Database connection string
	connStr := "host=db port=5432 user=latte password=latte dbname=frappuccino sslmode=disable"
	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	// Прокидываем зависимости в middleware
	mux := http.NewServeMux()

	// Обработка запросов на статичные файлы (HTML, CSS, JS)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/templates"))))

	// Главный маршрут, который будет открывать файл catalog.html
	mux.Handle("/", middleware.LoginOrLastVisitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "catalog.html"))
	})))
	// mux.HandleFunc("/chat", middleware.LoginOrLastVisitHandler)

	// Запуск сервера
	log.Println("Сервер работает на порту 8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
	}
}
