package services

import (
	"fmt"
	"log/slog"
	"time"

	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
)

// PostService provides operations for managing posts.
type PostService struct {
	PostRepository ports.PostRepository
	SessionRepo    ports.SessionRepository
}

// NewPostService creates a new instance of PostService.
func NewPostService(postRepo ports.PostRepository, sessionRepo ports.SessionRepository) *PostService {
	return &PostService{
		PostRepository: postRepo,
		SessionRepo:    sessionRepo,
	}
}

// CreatePost creates a new post using session data.
func (s *PostService) CreatePost(sessionID, title, text, imageURL string) (*models.Post, error) {
	userData, ok := s.SessionRepo.GetSessionData(sessionID)
	if !ok {
		slog.Warn("Session not found", "SessionID", sessionID)
		return nil, fmt.Errorf("session not found")
	}

	post := &models.Post{
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
		slog.Error("Failed to create post", "error", err)
		return nil, err
	}

	slog.Info("Post created successfully", "PostID", createdPost.ID, "UserName", post.UserName)
	return createdPost, nil
}

// GetAllPosts retrieves all posts.
func (s *PostService) GetAllPosts() ([]*models.Post, error) {
	posts, err := s.PostRepository.GetAllPosts()
	if err != nil {
		slog.Error("Failed to fetch posts", "error", err)
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}
	slog.Info("All posts retrieved", "Count", len(posts))
	return posts, nil
}

// GetPostByID retrieves a post by its ID.
func (s *PostService) GetPostByID(id string) (*models.Post, error) {
	post, err := s.PostRepository.GetPostByID(id)
	if err != nil {
		slog.Error("Failed to get post by ID", "PostID", id, "error", err)
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}
	slog.Info("Post retrieved", "PostID", id)
	return post, nil
}

// StartArchiver runs a background job that periodically checks posts for archiving.
func (s *PostService) StartArchiver() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C

			posts, err := s.PostRepository.GetAllPosts()
			if err != nil {
				slog.Error("Error fetching posts for archiving", "error", err)
				continue
			}

			now := time.Now()

			for _, post := range posts {
				if post.ArchivedAt != nil {
					continue
				}

				shouldArchive := false

				if len(post.Comments) == 0 {
					if post.CreatedAt.Add(10 * time.Minute).Before(now) {
						shouldArchive = true
					}
				} else {
					var lastCommentTime time.Time
					for _, c := range post.Comments {
						if c.CreatedAt.After(lastCommentTime) {
							lastCommentTime = c.CreatedAt
						}
					}
					if lastCommentTime.Add(15 * time.Minute).Before(now) {
						shouldArchive = true
					}
				}

				if shouldArchive {
					err := s.PostRepository.ArchivePost(post.ID)
					if err != nil {
						slog.Error("Failed to archive post", "PostID", post.ID, "error", err)
					} else {
						slog.Info("Post archived", "PostID", post.ID)
					}
				}
			}
		}
	}()
}

// GetArchivedPosts returns all archived posts.
func (s *PostService) GetArchivedPosts() ([]*models.Post, error) {
	posts, err := s.PostRepository.GetArchivedPosts()
	if err != nil {
		slog.Error("Failed to fetch archived posts", "error", err)
		return nil, err
	}
	slog.Info("Archived posts retrieved", "Count", len(posts))
	return posts, nil
}

// GetArchivedPostByID retrieves an archived post by its ID.
func (s *PostService) GetArchivedPostByID(id string) (*models.Post, error) {
	post, err := s.PostRepository.GetArchivedPostByID(id)
	if err != nil {
		slog.Error("Failed to get archived post by ID", "PostID", id, "error", err)
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}
	slog.Info("Archived post retrieved", "PostID", id)
	return post, nil
}
