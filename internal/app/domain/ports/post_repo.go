// /internal/app/domain/ports/post_repository.go
package ports

import "1337b04rd/internal/app/domain/models"

// Интерфейс репозитория для работы с постами
type PostRepository interface {
	CreatePost(post *models.Post) (*models.Post, error)
	UpdatePost(post *models.Post) (*models.Post, error)
	DeletePost(id string) error
	GetPostByID(id string) (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
}
