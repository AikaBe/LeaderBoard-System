package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
	"1337b04rd/internal/app/domain/services"
)

type CommentHandler struct {
	CommentService *services.CommentService
	GetSessionData ports.SessionRepository
	S3Adapter      ports.S3Adapter
}

func NewCommentHandler(commentService *services.CommentService, sessionRepo ports.SessionRepository, s3Adapter ports.S3Adapter) *CommentHandler {
	return &CommentHandler{
		CommentService: commentService,
		GetSessionData: sessionRepo,
		S3Adapter:      s3Adapter,
	}
}

type CreateCommentResponse struct {
	ID         int    `json:"id"`
	UserAvatar string `json:"userAvatar"`
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	// Get session ID from context
	sessionIDRaw := r.Context().Value("sessionId")
	sessionID, ok := sessionIDRaw.(string)
	if !ok || sessionID == "" {
		http.Error(w, "Unauthorized: no session ID", http.StatusUnauthorized)
		return
	}

	// Get user data
	userData, ok := h.GetSessionData.GetSessionData(sessionID)
	if !ok || userData.Name == "" {
		http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)
		return
	}

	// Parse form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values
	postIDStr := r.FormValue("post_id")
	text := r.FormValue("comment")
	parentIDStr := r.FormValue("parent_id")

	if postIDStr == "" || text == "" {
		http.Error(w, "Missing post ID or comment text", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var parentID *int
	if parentIDStr != "" {
		id, err := strconv.Atoi(parentIDStr)
		if err != nil {
			http.Error(w, "Invalid parent ID", http.StatusBadRequest)
			return
		}
		parentID = &id
	}

	// Upload image
	imageURL, err := h.S3Adapter.UploadImage(r, "comment")
	if err != nil {
		slog.Error("Image upload failed", "error", err)
		http.Error(w, "Image upload failed", http.StatusInternalServerError)
		return
	}

	// Create comment
	comment := models.Comment{
		PostID:          postID,
		ParentCommentID: parentID,
		UserName:        userData.Name,
		UserAvatar:      userData.Avatar,
		Text:            text,
		ImageURL:        imageURL,
		CreatedAt:       time.Now(),
	}

	if _, err := h.CommentService.CreateComment(comment); err != nil {
		slog.Error("Failed to create comment", "error", err)
		http.Error(w, "Failed to save comment", http.StatusInternalServerError)
		return
	}

	slog.Info("Comment created successfully", "postID", postID, "user", userData.Name)

	// Redirect to post
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}
