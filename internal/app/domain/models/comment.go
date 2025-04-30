package models

import (
	"time"
)

type Comment struct {
	ID         int       `json:"id"`         // Уникальный ID комментария
	PostID     string    `json:"postID"`     // ID поста, к которому относится комментарий
	ParentID   *string   `json:"parentID"`   // ID родительского комментария (если nil — это ответ на пост)
	UserName   string    `json:"userName"`   // Имя пользователя
	UserAvatar string    `json:"userAvatar"` // Аватар пользователя
	Text       string    `json:"text"`       // Текст комментария
	CreatedAt  time.Time `json:"createdAt"`  // Дата и время создания
}
