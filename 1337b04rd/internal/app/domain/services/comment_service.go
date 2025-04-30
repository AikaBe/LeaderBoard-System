package services

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
	"errors"
	"time"
)

// Структура для CommentService
type CommentService struct {
	CommentRepo ports.CommentRepository
}

// Конструктор для CommentService
func NewCommentService(repo ports.CommentRepository) *CommentService {
	return &CommentService{
		CommentRepo: repo,
	}
}

// Метод для создания комментария

func (s *CommentService) CreateComment(comment models.Comment) (*models.Comment, error) {
	if comment.PostID == "" || comment.UserName == "" || comment.Text == "" {
		return nil, errors.New("missing required fields (PostID, UserName, Text)")
	}

	// Установим дату, если не пришла
	if comment.CreatedAt.IsZero() {
		comment.CreatedAt = time.Now()
	}

	return s.CommentRepo.CreateComment(comment)
}

func (s *CommentService) GetCommentsByPostID(postID string) ([]*models.Comment, error) {
	if postID == "" {
		return nil, errors.New("invalid PostID")
	}
	return s.CommentRepo.GetCommentsByPostID(postID)
}

func (s *CommentService) DeleteComment(commentID string) error {
	if commentID == "" {
		return errors.New("invalid CommentID")
	}
	return s.CommentRepo.DeleteComment(commentID)
}
