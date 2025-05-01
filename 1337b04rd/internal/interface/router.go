package router

import (
	"1337b04rd/internal/adapters/database"
	"1337b04rd/internal/app/domain/services"
	"1337b04rd/internal/interface/handlers"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)

// SetupRoutes конфигурирует маршруты для HTTP-сервера
func SetupRoutes() {
	// Создаем подключение к базе данных
	db, err := sql.Open("postgres", "user=username password=password dbname=yourdb sslmode=disable")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	// Проверка подключения к базе данных
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	// Создаем репозиторий
	postRepo := database.NewPostRepositoryPg(db)

	// Создаем сервис для работы с постами
	postService := services.NewPostService(postRepo, nil) // Замените nil на репозиторий сессий, если нужно

	// Создаем хэндлер для обработки запросов
	postHandler := handlers.NewPostHandler(postService)

	// Настройка маршрутов
	http.HandleFunc("/submit-post", postHandler.CreatePost)
	// Здесь можно добавить другие маршруты
	http.HandleFunc("/catalog", postHandler.GetAllPosts)

	// Запуск сервера
	log.Fatal(http.ListenAndServe(":8080", nil))
}
