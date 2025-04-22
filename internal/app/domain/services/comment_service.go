package services

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
)

type CommentService struct {
	CommentRepo ports.CommentRepository
}

func NewCommentService(repo ports.CommentRepository) *CommentService {
	return &CommentService{
		CommentRepo: repo,
	}
}

func (s *CommentService) CreateComment(comment models.Comment) (*models.Comment, error) {
	// здесь можно добавить любую логику валидации, если нужно
	return s.CommentRepo.CreateComment(comment)
}

func (s *CommentService) GetCommentsByPostID(postID string) ([]*models.Comment, error) {
	return s.CommentRepo.GetCommentsByPostID(postID)
}

func (s *CommentService) DeleteComment(commentID string) error {
	return s.CommentRepo.DeleteComment(commentID)
}
