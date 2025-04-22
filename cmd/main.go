package main

import (
	d "1337b04rd/internal/adapters/database"
	"1337b04rd/internal/app/domain/ports"
	"1337b04rd/internal/app/domain/services"
	"1337b04rd/internal/interface/handlers"
	"1337b04rd/internal/interface/middleware"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

func initRepository(dsn string) (ports.PostRepository, *sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to database: %v", err)
	}

	postRepo := d.NewPostRepositoryPg(db)

	return postRepo, db, nil
}

func main() {
	// DSN строка подключения (пример)
	dsn := "host=localhost port=5432 user=latte password=latte dbname=frappuccino sslmode=disable"

	postRepo, db, err := initRepository(dsn)
	if err != nil {
		log.Fatalf("Ошибка инициализации репозитория: %v", err)
	}
	defer db.Close()

	// Создаем сервис для работы с постами
	postService := services.NewPostService(postRepo)

	// Создаем обработчик
	postHandler := handlers.NewPostHandler(postService)

	// Прокидываем зависимости в middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/login", middleware.LoginHandler)
	mux.HandleFunc("/last-visit", middleware.LastVisitHandler)

	// Привязываем обработчики к маршрутам
	mux.HandleFunc("/posts", postHandler.GetAllPosts)      // Получение всех постов
	mux.HandleFunc("/create-post", postHandler.CreatePost) // Создание нового поста

	// Запуск сервера
	log.Println("Сервер работает на порту 8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
	}
}
