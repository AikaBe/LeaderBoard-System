package services

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
	"errors"
	"log/slog"
	"strconv"
	"time"
)

// CommentService provides methods to work with comments.
type CommentService struct {
	CommentRepo ports.CommentRepository
}

// NewCommentService creates a new instance of CommentService.
func NewCommentService(repo ports.CommentRepository) *CommentService {
	return &CommentService{
		CommentRepo: repo,
	}
}

// CreateComment validates and creates a new comment.
func (s *CommentService) CreateComment(comment models.Comment) (*models.Comment, error) {
	if comment.PostID == 0 || comment.UserName == "" || comment.Text == "" {
		slog.Warn("Missing required fields in comment creation", "PostID", comment.PostID, "UserName", comment.UserName, "Text", comment.Text)
		return nil, errors.New("missing required fields (PostID, UserName, Text)")
	}

	if comment.CreatedAt.IsZero() {
		comment.CreatedAt = time.Now()
	}

	createdComment, err := s.CommentRepo.CreateComment(comment)
	if err != nil {
		slog.Error("Failed to create comment", "error", err)
		return nil, err
	}

	slog.Info("Comment created successfully", "CommentID", createdComment.ID)
	return createdComment, nil
}

// GetCommentsByPostID returns all comments associated with a specific post ID.
func (s *CommentService) GetCommentsByPostID(postID int) ([]*models.Comment, error) {
	if postID == 0 {
		slog.Warn("Invalid post ID provided for comment retrieval")
		return nil, errors.New("invalid PostID")
	}

	comments, err := s.CommentRepo.GetCommentsByPostID(postID)
	if err != nil {
		slog.Error("Failed to retrieve comments", "PostID", postID, "error", err)
		return nil, err
	}

	slog.Info("Comments retrieved successfully", "PostID", postID, "Count", len(comments))
	return comments, nil
}

// DeleteComment deletes a comment by its ID.
func (s *CommentService) DeleteComment(commentID int) error {
	if commentID == 0 {
		slog.Warn("Invalid comment ID provided for deletion")
		return errors.New("invalid CommentID")
	}

	commentIDString := strconv.Itoa(commentID)
	err := s.CommentRepo.DeleteComment(commentIDString)
	if err != nil {
		slog.Error("Failed to delete comment", "CommentID", commentID, "error", err)
		return err
	}

	slog.Info("Comment deleted successfully", "CommentID", commentID)
	return nil
}
