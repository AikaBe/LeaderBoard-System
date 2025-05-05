package ports

import "1337b04rd/internal/app/domain/models"

type CommentRepository interface {
	CreateComment(comment models.Comment) (*models.Comment, error)
	GetCommentsByPostID(postID int) ([]*models.Comment, error)
	DeleteComment(commentID string) error
}
