package domain

import (
	"1337b04rd/internal/app/domain/models"
	"errors"
	"testing"
	"time"
)

type MockCommentRepository struct {
	comments     map[int]*models.Comment
	lastID       int
	createErr    error
	getErr       error
	deleteErr    error
	updateErr    error
	listErr      error
	setHiddenErr error
}

func NewMockCommentRepository() *MockCommentRepository {
	return &MockCommentRepository{
		comments: make(map[int]*models.Comment),
		lastID:   0,
	}
}

func (m *MockCommentRepository) CreateComment(comment *models.Comment) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.lastID++
	comment.ID = m.lastID
	comment.CreatedAt = time.Now()
	m.comments[comment.ID] = comment
	return nil
}

func (m *MockCommentRepository) GetCommentByID(id int) (*models.Comment, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	comment, exists := m.comments[id]
	if !exists {
		return nil, errors.New("comment not found")
	}
	return comment, nil
}

func (m *MockCommentRepository) DeleteComment(id int) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, exists := m.comments[id]; !exists {
		return errors.New("comment not found")
	}
	delete(m.comments, id)
	return nil
}

func (m *MockCommentRepository) GetCommentsByPostID(postID int) ([]*models.Comment, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var postComments []*models.Comment
	for _, c := range m.comments {
		if c.PostID == postID {
			postComments = append(postComments, c)
		}
	}
	return postComments, nil
}

func TestCreateComment(t *testing.T) {
	repo := NewMockCommentRepository()
	comment := &models.Comment{PostID: 1, Text: "Test"}

	err := repo.CreateComment(comment)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if comment.ID != 1 {
		t.Errorf("expected ID to be 1, got %d", comment.ID)
	}
	if comment.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}
