package models

import (
	"time"
)

// models/post.go
type Post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
	UserName   string    `json:"user_name"`
	UserAvatar string    `json:"user_avatar"`
	ImageURL   string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ArchivedAt *time.Time
	IsHidden   bool `json:"is_hidden"`
	Comments   []*Comment
}
