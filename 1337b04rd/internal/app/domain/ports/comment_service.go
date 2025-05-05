package ports

import "1337b04rd/internal/app/domain/models"

type CommentService interface {
	GetCommentsByPostID(postID int) ([]*models.Comment, error)
}
