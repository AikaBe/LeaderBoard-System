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

type Comment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
