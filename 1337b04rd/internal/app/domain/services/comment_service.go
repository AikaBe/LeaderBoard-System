package services

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
	"errors"
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
	// Пример базовой валидации (можно расширить)
	if comment.PostID == "" || comment.UserName == "" || comment.Text == "" {
		return nil, errors.New("missing required fields (PostID, UserName, Text)")
	}

	// Вызываем репозиторий для создания комментария
	return s.CommentRepo.CreateComment(comment)
}

// Метод для получения комментариев по ID поста
func (s *CommentService) GetCommentsByPostID(postID string) ([]*models.Comment, error) {
	// Пример дополнительной проверки на пустой postID
	if postID == "" {
		return nil, errors.New("invalid PostID")
	}
	return s.CommentRepo.GetCommentsByPostID(postID)
}

// Метод для удаления комментария
func (s *CommentService) DeleteComment(commentID string) error {
	// Пример проверки на пустой commentID
	if commentID == "" {
		return errors.New("invalid CommentID")
	}
	return s.CommentRepo.DeleteComment(commentID)
}
