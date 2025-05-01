package main

import (
	"1337b04rd/internal/interface/middleware"
	"database/sql"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection string
	connStr := "host=db port=5432 user=board_user password=board_pass dbname=board_db sslmode=disable"

	// Подключение к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error trying to connect to database: ", err)
	}

	// HTTP роутинг
	mux := http.NewServeMux()

	// Обработка запросов на статичные файлы (HTML, CSS, JS)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/templates"))))

	// Главный маршрут, который будет открывать файл catalog.html
	mux.Handle("/", middleware.LoginOrLastVisitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "catalog.html"))
	})))
	// mux.HandleFunc("/chat", middleware.LoginOrLastVisitHandler)
	mux.Handle("/create-post", middleware.LoginOrLastVisitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "create-post.html"))
	})))
	mux.Handle("/archive", middleware.LoginOrLastVisitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "archive.html"))
	})))
	mux.Handle("/archive-post", middleware.LoginOrLastVisitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "archive-post.html"))
	})))

	mux.Handle("/post", middleware.LoginOrLastVisitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "post.html"))
	})))
	mux.Handle("/error", middleware.LoginOrLastVisitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "error.html"))
	})))

	log.Println("The server is running on port 8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
