package handlers

import (
	"1337b04rd/internal/app/domain/ports"
	"1337b04rd/internal/app/domain/services"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// PostHandler представляет структуру для работы с HTTP-запросами для постов
type PostHandler struct {
	PostService *services.PostService
	S3Adapter   ports.S3Adapter
}

// Новый конструктор для создания нового PostHandler
func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		PostService: postService,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Сначала извлечь поля формы
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("subject")
	text := r.FormValue("comment")
	sessionID := r.FormValue("name")

	// Загрузка картинки — вся логика внутри адаптера
	imageURL, err := h.S3Adapter.UploadImage(r)
	if err != nil {
		http.Error(w, "Image upload failed", http.StatusInternalServerError)
		return
	}

	// Создание поста
	createdPost, err := h.PostService.CreatePost(sessionID, title, text, imageURL)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPost)
}

// Метод для обработки обновления имени пользователя в посте
func (h *PostHandler) UpdateUserName(w http.ResponseWriter, r *http.Request) {
	// Получаем ID поста из URL
	postID := r.URL.Query().Get("postID")
	if postID == "" {
		http.Error(w, "Missing postID in URL", http.StatusBadRequest)
		return
	}

	var request struct {
		NewUserName string `json:"newUserName"`
	}

	// Декодируем новое имя пользователя из тела запроса
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обновляем имя пользователя в посте и комментариях
	err = h.PostService.UpdateUserName(postID, request.NewUserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User name updated successfully"))
}

// Метод для обработки удаления поста
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	// Получаем ID поста из URL
	postID := r.URL.Query().Get("postID")
	if postID == "" {
		http.Error(w, "Missing postID in URL", http.StatusBadRequest)
		return
	}

	// Удаляем пост через сервис
	err := h.PostService.DeletePost(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post deleted successfully"))
}

// Метод для обработки запроса получения всех постов
func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostService.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("1337b04rd/web/templates/catalog.html") // путь к твоему HTML-шаблону
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}

// Запланированное удаление поста
func (h *PostHandler) schedulePostDeletion(postID string, delay time.Duration) {
	time.Sleep(delay)
	// Делаем попытку удалить пост после задержки
	err := h.PostService.DeletePost(postID)
	if err != nil {
		// Логируем ошибку, если не удалось удалить пост
		fmt.Println("Failed to delete post:", err)
	}
}
