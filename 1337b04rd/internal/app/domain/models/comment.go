package models

import (
	"time"
)

type Comment struct {
	ID              int        `json:"id"`
	PostID          int        `json:"post_id"`
	ParentCommentID *int       `json:"parent_comment_id"` // Может быть NULL
	UserName        string     `json:"user_name"`
	UserAvatar      string     `json:"user_avatar"`
	Text            string     `json:"text"`
	ImageURL        string     `json:"image_url,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	Replies         []*Comment `json:"replies"`
}
