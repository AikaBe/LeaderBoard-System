package main

import (
	d "1337b04rd/internal/adapters/database"
	"1337b04rd/internal/app/domain/ports"
	"1337b04rd/internal/app/domain/services"
	"1337b04rd/internal/interface/handlers"
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
	// DSN строка подключения
	dsn := "host=db port=5432 user=latte password=latte dbname=frappuccino sslmode=disable"

	postRepo, commentRepo, db, err := initRepository(dsn)
	if err != nil {
		log.Fatalf("Ошибка инициализации репозитория: %v", err)
	}
	defer db.Close()

	// Сервисы
	postService := services.NewPostService(postRepo, nil)
	commentService := services.NewCommentService(commentRepo)

	// Хэндлеры
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)

	// HTTP роутинг
	mux := http.NewServeMux()

	// Роуты для постов
	mux.HandleFunc("/posts", postHandler.GetAllPosts)      // Получение всех постов
	mux.HandleFunc("/create-post", postHandler.CreatePost) // Создание нового поста

	// Роуты для комментариев
	mux.HandleFunc("/comments", commentHandler.CreateComment) // Создание нового комментария

	// Запуск сервера
	log.Println("Сервер работает на порту 8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Ошибка при запуске сервера: ", err)
	}
}

// package main

// import (
// 	"1337b04rd/internal/interface/middleware"
// 	"database/sql"
// 	"log"
// 	"net/http"
// 	"path/filepath"

// 	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
// )

// func main() {
// 	// Database connection string
// 	connStr := "host=db port=5432 user=latte password=latte dbname=frappuccino sslmode=disable"
// 	// Open a connection to the database
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Error connecting to the database: ", err)
// 	}
// 	defer db.Close()

// 	if err := db.Ping(); err != nil {
// 		log.Fatal("Error pinging the database: ", err)
// 	}

// 	// Прокидываем зависимости в middleware
// 	mux := http.NewServeMux()

// 	// Обработка запросов на статичные файлы (HTML, CSS, JS)
// 	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/templates"))))

// 	// Главный маршрут, который будет открывать файл catalog.html
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		http.ServeFile(w, r, filepath.Join("web", "templates", "catalog.html"))
// 	})
// 	mux.HandleFunc("/chat", middleware.LoginOrLastVisitHandler)

// 	// Запуск сервера
// 	log.Println("Сервер работает на порту 8080...")
// 	err = http.ListenAndServe(":8080", mux)
// 	if err != nil {
// 		log.Fatal("Ошибка при запуске сервера: ", err)
// 	}
// }
