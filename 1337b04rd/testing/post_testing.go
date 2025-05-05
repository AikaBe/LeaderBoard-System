package testing

import (
	"1337b04rd/internal/app/domain/models"
	"errors"
	"testing"
	"time"
)

type MockPostRepository struct {
	posts          map[int]*models.Post
	lastID         int
	createErr      error
	getErr         error
	deleteErr      error
	listErr        error
	updateErr      error
	setArchivedErr error
}

func NewMockPostRepository() *MockPostRepository {
	return &MockPostRepository{
		posts:  make(map[int]*models.Post),
		lastID: 0,
	}
}

func (m *MockPostRepository) CreatePost(post *models.Post) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.lastID++
	post.ID = m.lastID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) GetPostByID(id int) (*models.Post, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	post, exists := m.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (m *MockPostRepository) DeletePost(id int) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, exists := m.posts[id]; !exists {
		return errors.New("post not found")
	}
	delete(m.posts, id)
	return nil
}

// Helper function to create a valid post
func createValidPost() *models.Post {
	return &models.Post{
		Title:      "Test Post",
		Text:       "This is a test post content",
		UserName:   "Test User",
		UserAvatar: "test-avatar-url",
		ImageURL:   "test-image-url",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsHidden:   false,
	}
}

// Test functions

func TestCreatePost_Success(t *testing.T) {
	// Setup
	mockRepo := NewMockPostRepository()
	testPost := createValidPost()

	// Execute
	err := mockRepo.CreatePost(testPost)
	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if testPost.ID == 0 {
		t.Error("Expected post ID to be set, but it wasn't")
	}
	if storedPost, _ := mockRepo.GetPostByID(testPost.ID); storedPost == nil {
		t.Error("Post was not stored in repository")
	}
}

func TestCreatePost_ValidationErrors(t *testing.T) {
	// Setup
	mockRepo := NewMockPostRepository()

	// Test cases for validation errors
	testCases := []struct {
		name        string
		modifyPost  func(*models.Post)
		expectedErr string
	}{
		{
			name: "Missing title",
			modifyPost: func(p *models.Post) {
				p.Title = ""
			},
			expectedErr: "title is required",
		},
		{
			name: "Missing text",
			modifyPost: func(p *models.Post) {
				p.Text = ""
			},
			expectedErr: "text is required",
		},
		{
			name: "Missing user name",
			modifyPost: func(p *models.Post) {
				p.UserName = ""
			},
			expectedErr: "user name is required",
		},
		{
			name: "Missing avatar",
			modifyPost: func(p *models.Post) {
				p.UserAvatar = ""
			},
			expectedErr: "user avatar is required",
		},
		{
			name: "Missing image URL",
			modifyPost: func(p *models.Post) {
				p.ImageURL = ""
			},
			expectedErr: "image URL is required",
		},
	}

	// Execute and verify each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			postCopy := createValidPost()
			tc.modifyPost(postCopy)
			err := mockRepo.CreatePost(postCopy)

			if err == nil {
				t.Fatalf("Expected error but got nil")
			}
			if err.Error() != tc.expectedErr {
				t.Fatalf("Expected error '%s', got '%s'", tc.expectedErr, err.Error())
			}
		})
	}
}

func TestGetPostByID_Success(t *testing.T) {
	// Setup
	mockRepo := NewMockPostRepository()
	testPost := createValidPost()
	_ = mockRepo.CreatePost(testPost)

	// Execute
	retrievedPost, err := mockRepo.GetPostByID(testPost.ID)
	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if retrievedPost == nil {
		t.Fatal("Expected post to be returned, got nil")
	}
	if retrievedPost.ID != testPost.ID {
		t.Fatalf("Expected post ID %d, got %d", testPost.ID, retrievedPost.ID)
	}
	if retrievedPost.Title != testPost.Title {
		t.Fatalf("Expected title '%s', got '%s'", testPost.Title, retrievedPost.Title)
	}
}

func TestDeletePost_Success(t *testing.T) {
	// Setup
	mockRepo := NewMockPostRepository()
	testPost := createValidPost()
	_ = mockRepo.CreatePost(testPost)

	// Execute
	err := mockRepo.DeletePost(testPost.ID)
	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify post is deleted
	_, getErr := mockRepo.GetPostByID(testPost.ID)
	if getErr == nil {
		t.Fatal("Expected post to be deleted, but it wasn't")
	}
}
