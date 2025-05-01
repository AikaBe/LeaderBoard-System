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
	db, err := sql.Open("postgres", "user=username password=password dbname=yourdb sslmode=disable")
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Ping DB error:", err)
	}

	postRepo := database.NewPostRepositoryPg(db)
	postService := services.NewPostService(postRepo)
	s3Adapter := s3.NewS3Adapter("./images")

	postHandler := handlers.NewPostHandler(postService, s3Adapter)

	http.HandleFunc("/submit-post", postHandler.CreatePost)
	log.Println("Server listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
