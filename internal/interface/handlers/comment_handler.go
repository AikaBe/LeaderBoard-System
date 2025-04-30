package handlers

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/services"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Структура CommentHandler
type CommentHandler struct {
	CommentService *services.CommentService
}

// Конструктор для CommentHandler
func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
	}
}

// Структура ответа для создания комментария
type CreateCommentResponse struct {
	ID         int    `json:"id"`         // ID нового комментария
	UserAvatar string `json:"userAvatar"` // Аватар пользователя
}

// Метод для создания комментария
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid comment format", http.StatusBadRequest)
		return
	}

	// Проверяем, что дата корректно парсится
	if comment.CreatedAt.IsZero() {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	// Если клиент не передал дату — устанавливаем её
	if comment.CreatedAt.IsZero() {
		comment.CreatedAt = time.Now()
	}

	createdComment, err := h.CommentService.CreateComment(comment)
	if err != nil {
		log.Printf("error creating comment: %v", err) // ← добавь это
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	// Формируем ответ с нужными полями
	response := CreateCommentResponse{
		ID:         createdComment.ID,
		UserAvatar: createdComment.UserAvatar,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
