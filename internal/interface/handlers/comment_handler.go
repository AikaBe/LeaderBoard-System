package handlers

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/services"
	"encoding/json"
	"net/http"
)

// Структура CommentHandler
type CommentHandler struct {
	CommentService *services.CommentService
}

// Конструктор (фабрика) для CommentHandler
func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
	}
}

// Метод CreateComment
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid comment format", http.StatusBadRequest)
		return
	}

	createdComment, err := h.CommentService.CreateComment(comment)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdComment)
}
