package models

import (
	"time"
)

// models/post.go
type Post struct {
	ID         string
	Title      string
	Text       string
	UserID     string
	UserName   string
	UserAvatar string
	ImageURL   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsHidden   bool // <--- добавляем это поле
	Comments   []Comment
}
