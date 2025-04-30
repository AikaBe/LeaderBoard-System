package database

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
	"database/sql"
	"fmt"
	"strconv"
)

type CommentRepositoryPg struct {
	db *sql.DB
}

func NewCommentRepositoryPg(db *sql.DB) ports.CommentRepository {
	return &CommentRepositoryPg{db: db}
}

func (r *CommentRepositoryPg) CreateComment(comment models.Comment) (*models.Comment, error) {
	// 1. Проверка: существует ли пост
	postIDInt, err := strconv.Atoi(comment.PostID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID format: %v", err)
	}

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`
	err = r.db.QueryRow(query, postIDInt).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking post existence: %v", err)
	}

	if !exists {
		return nil, fmt.Errorf("post with id %d does not exist", postIDInt)
	}

	// 2. Вставка комментария
	query = `INSERT INTO comments (post_id, parent_comment_id, user_name, user_avatar, text, created_at)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int
	err = r.db.QueryRow(
		query,
		postIDInt,
		comment.ParentID,
		comment.UserName,
		comment.UserAvatar,
		comment.Text,
		comment.CreatedAt,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error creating comment: %v", err)
	}

	comment.ID = id
	return &comment, nil
}

// Метод для получения комментариев по ID поста
func (r *CommentRepositoryPg) GetCommentsByPostID(postID string) ([]*models.Comment, error) {
	// Конвертируем postID в int
	idInt, err := strconv.Atoi(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid post id: %v", err)
	}

	query := `SELECT id, post_id, parent_comment_id, user_name, user_avatar, text, created_at
              FROM comments
              WHERE post_id = $1
              ORDER BY created_at ASC`

	rows, err := r.db.Query(query, idInt)
	if err != nil {
		return nil, fmt.Errorf("error getting comments: %v", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.ParentID,
			&c.UserName,
			&c.UserAvatar,
			&c.Text,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %v", err)
		}
		comments = append(comments, &c)
	}

	return comments, nil
}

// Метод для удаления комментария
func (r *CommentRepositoryPg) DeleteComment(commentID string) error {
	idInt, err := strconv.Atoi(commentID)
	if err != nil {
		return fmt.Errorf("invalid comment id: %v", err)
	}

	query := `DELETE FROM comments WHERE id = $1`
	_, err = r.db.Exec(query, idInt)
	if err != nil {
		return fmt.Errorf("error deleting comment: %v", err)
	}

	return nil
}
