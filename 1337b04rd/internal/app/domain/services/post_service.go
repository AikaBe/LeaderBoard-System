// /internal/app/services/post_service.go
package services

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
	"fmt"
	"time"
)

// Сервис для работы с постами
type PostService struct {
	PostRepository ports.PostRepository // Используем интерфейс репозитория
	SessionRepo    ports.SessionRepository
}

// Конструктор для создания нового PostService
func NewPostService(postRepo ports.PostRepository, sessionRepo ports.SessionRepository) *PostService {
	return &PostService{
		PostRepository: postRepo,
		SessionRepo:    sessionRepo,
	}
}

// Метод для проверки, истекло ли время с последнего комментария поста
func (s *PostService) IsPostExpired(p *models.Post) bool {
	if len(p.Comments) == 0 {
		// Если прошло больше 10 минут с создания, пост истек
		return time.Since(p.CreatedAt) > 10*time.Minute
	}
	// Если прошло больше 15 минут с последнего комментария, пост истек
	lastComment := p.Comments[len(p.Comments)-1]
	return time.Since(lastComment.CreatedAt) > 15*time.Minute
}

func (s *PostService) CreatePost(sessionID, title, text, imageURL string) (*models.Post, error) {
	userData, ok := s.SessionRepo.GetSessionData(sessionID)
	if !ok {
		return nil, fmt.Errorf("session not found")
	}

	postID := generateUniqueID()
	post := &models.Post{
		ID:         postID,
		Title:      title,
		Text:       text,
		UserName:   userData.Name,
		UserAvatar: userData.Avatar,
		ImageURL:   imageURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	createdPost, err := s.PostRepository.CreatePost(post)
	if err != nil {
		return nil, err
	}

	delay := 10 * time.Minute
	if len(post.Comments) > 0 {
		delay = 15 * time.Minute
	}
	go s.schedulePostDeletion(post.ID, delay)

	return createdPost, nil
}

// Обновление имени пользователя в посте и комментариях
func (s *PostService) UpdateUserName(postID, newUserName string) error {
	post, err := s.PostRepository.GetPostByID(postID)
	if err != nil {
		return err
	}

	// Обновляем имя пользователя в посте и комментариях
	post.UserName = newUserName
	for i := range post.Comments {
		post.Comments[i].UserName = newUserName
	}

	_, err = s.PostRepository.UpdatePost(post)
	return err
}

// Метод для удаления поста после определенного времени
func (s *PostService) schedulePostDeletion(postID string, delay time.Duration) {
	time.Sleep(delay)

	post, err := s.PostRepository.GetPostByID(postID)
	if err != nil || post == nil {
		return
	}

	if s.IsPostExpired(post) {
		post.IsHidden = true
		post.UpdatedAt = time.Now()
		_, _ = s.PostRepository.UpdatePost(post) // обновляем флаг скрытия
	}
}

// Метод для удаления поста
func (s *PostService) DeletePost(postID string) error {
	// Проверяем, существует ли пост
	post, err := s.PostRepository.GetPostByID(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return fmt.Errorf("post not found")
	}

	// Удаляем пост из репозитория
	err = s.PostRepository.DeletePost(postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	return nil
}

// Метод для получения всех постов
func (s *PostService) GetAllPosts() ([]*models.Post, error) {
	posts, err := s.PostRepository.GetAllPosts()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}
	return posts, nil
}

// Генерация уникального ID для поста
func generateUniqueID() string {
	return fmt.Sprintf("post-%d", time.Now().UnixNano())
}
