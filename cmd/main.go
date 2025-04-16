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

// // Строка подключения к PostgreSQL
// connStr := "postgres://user:password@localhost:5432/mydb?sslmode=disable"
// db, err := sql.Open("postgres", connStr)
// if err != nil {
// 	log.Fatal("Unable to connect to database:", err)
// }
// defer db.Close()

// // Проверка подключения
// err = db.Ping()
// if err != nil {
// 	log.Fatal("Unable to ping database:", err)
// }

// // Создаем репозиторий
// postRepo := repositories.NewPostRepositoryPg(db)

// // Создаем сервис с репозиторием
// postService := services.NewPostService(postRepo)

// // Пример работы с сервисом
// _, err = postService.CreatePost("Post title", "Post text", "user123", "User Name", "User Avatar", "http://example.com/image.jpg")
// if err != nil {
// 	log.Fatal("Error creating post:", err)
// }

// fmt.Println("Post created successfully!")
// }
