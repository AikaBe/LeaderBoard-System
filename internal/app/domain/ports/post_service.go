// /internal/app/domain/ports/post_service.go
package ports

import "1337b04rd/internal/app/domain/models"

// Интерфейс сервиса для работы с постами
type PostService interface {
	CreatePost(title, text, userID string, userName, userAvatar, imageURL string) (*models.Post, error)
	UpdateUserName(postID, newUserName string) error
	IsPostExpired(post *models.Post) bool
}
