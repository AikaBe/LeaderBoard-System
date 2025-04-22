package models

import (
	"time"
)

type Comment struct {
	ID        string    `json:"id"`        // Уникальный ID комментария
	PostID    string    `json:"postID"`    // ID поста, к которому относится комментарий
	ParentID  *string   `json:"parentID"`  // Если это ответ на другой комментарий — его ID (может быть nil)
	UserID    string    `json:"userID"`    // ID пользователя
	UserName  string    `json:"userName"`  // Имя пользователя
	Text      string    `json:"text"`      // Текст комментария
	CreatedAt time.Time `json:"createdAt"` // Дата и время создания
}
