package handlers

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/services"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PostHandler представляет структуру для работы с HTTP-запросами для постов
type PostHandler struct {
	PostService *services.PostService
}

// Новый конструктор для создания нового PostHandler
func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		PostService: postService,
	}
}

// Метод для обработки создания поста
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	// Декодируем данные из тела запроса в структуру post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создаем новый пост через сервис
	createdPost, err := h.PostService.CreatePost(post.Title, post.Text, post.UserName, post.UserAvatar, post.ImageURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем созданный пост в формате JSON
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем список постов в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
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
