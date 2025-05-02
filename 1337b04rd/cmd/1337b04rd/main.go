package main

import (
	"1337b04rd/internal/adapters/database"
	d "1337b04rd/internal/adapters/database"
	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/app/domain/ports"
	"1337b04rd/internal/app/domain/services"
	"1337b04rd/internal/interface/handlers"
	"1337b04rd/internal/interface/middleware"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

func initRepository(dsn string) (ports.PostRepository, ports.CommentRepository, *sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error connecting to database: %v", err)
	}

	postRepo := d.NewPostRepositoryPg(db)
	commentRepo := d.NewCommentRepositoryPg(db)

	return postRepo, commentRepo, db, nil
}

func main() {
	// Database DSN
	dsn := "host=db port=5432 user=board_user password=board_pass dbname=board_db sslmode=disable"

	postRepo, commentRepo, db, err := initRepository(dsn)
	if err != nil {
		log.Fatalf("Ошибка инициализации репозитория: %v", err)
	}
	defer db.Close()

	// Создание необходимых сервисов

	commentService := services.NewCommentService(commentRepo)
	sessionRepo := &database.PostgresSessionRepo{DB: db}
	sessionService := services.NewSessionService(sessionRepo)
	postService := services.NewPostService(postRepo, sessionRepo)

	s3Adapter := &s3.Adapter{
		TripleSBaseURL: "http://triple-s:9000",
	}

	// Передаём адаптер в хэндлер
	postHandler := handlers.NewPostHandler(postService, s3Adapter)
	commentHandler := handlers.NewCommentHandler(commentService)

	// Роутинг
	mux := http.NewServeMux()

	// Статические файлы
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/templates"))))

	authMiddleware := middleware.AuthMiddleware{SessionService: sessionService}
	// Главная страница
	mux.Handle("/", authMiddleware.LoginOrLastVisitHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("web", "templates", "catalog.html"))
	})))

	// Посты
	mux.HandleFunc("/posts", postHandler.GetAllPosts)
	mux.HandleFunc("/create-post", postHandler.CreatePost)

	// Комментарии
	mux.HandleFunc("/comments", commentHandler.CreateComment)

	// Запуск сервера
	log.Println("Сервер работает на порту 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
