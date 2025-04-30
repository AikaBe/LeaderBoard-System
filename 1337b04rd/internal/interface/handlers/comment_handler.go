package handlers

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/services"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type CommentHandler struct {
	CommentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
	}
}

type CreateCommentResponse struct {
	ID         int    `json:"id"`
	UserAvatar string `json:"userAvatar"`
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	// Извлекаем заголовок с пользовательскими данными
	userHeader := r.Header.Get("X-User-Data")
	if userHeader == "" {
		http.Error(w, "Missing user data", http.StatusUnauthorized)
		return
	}

	// Парсим JSON из заголовка
	var userData models.UserData
	err := json.Unmarshal([]byte(userHeader), &userData)
	if err != nil || userData.Name == "" {
		http.Error(w, "Invalid user data", http.StatusUnauthorized)
		return
	}

	// Декодим тело запроса (только текст и PostID от клиента)
	var comment models.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil || comment.Text == "" || comment.PostID == "" {
		http.Error(w, "Invalid comment data", http.StatusBadRequest)
		return
	}

	// Заполняем поля, доверяя только userHeader
	comment.UserName = userData.Name
	comment.UserAvatar = userData.Avatar
	comment.CreatedAt = time.Now()

	createdComment, err := h.CommentService.CreateComment(comment)
	if err != nil {
		log.Printf("error creating comment: %v", err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	response := CreateCommentResponse{
		ID:         createdComment.ID,
		UserAvatar: createdComment.UserAvatar,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
