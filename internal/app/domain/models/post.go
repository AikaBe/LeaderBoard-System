package models

import (
	"time"
)

type Post struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Text       string    `json:"text"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name"`
	UserAvatar string    `json:"user_avatar"`
	ImageURL   string    `json:"image_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Comments   []Comment `json:"comments,omitempty"`
}

type Comment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
