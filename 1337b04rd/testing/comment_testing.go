package core

import (
	"errors"
	"time"
	"your_project/models"
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

func (m *MockCommentRepository) UpdateCommentText(commentID int, text string) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	comment, exists := m.comments[commentID]
	if !exists {
		return errors.New("comment not found")
	}
	comment.Text = text
	comment.UpdatedAt = time.Now()
	return nil
}

func (m *MockCommentRepository) SetHidden(commentID int, isHidden bool) error {
	if m.setHiddenErr != nil {
		return m.setHiddenErr
	}
	comment, exists := m.comments[commentID]
	if !exists {
		return errors.New("comment not found")
	}
	comment.IsHidden = isHidden
	return nil
}
