package database

import (
	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
)

// Реализация репозитория постов для PostgreSQL
type PostRepositoryPg struct {
	db *sql.DB // Подключение к базе данных
}

// Конструктор для создания нового PostRepositoryPg
func NewPostRepositoryPg(db *sql.DB) ports.PostRepository {
	return &PostRepositoryPg{db: db}
}

// Создание нового поста
func (r *PostRepositoryPg) CreatePost(post *models.Post) (*models.Post, error) {
	// Запрос на создание поста
	query := `INSERT INTO posts (id, title, text, user_id, user_name, user_avatar, image_url, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	// Выполняем запрос
	err := r.db.QueryRow(query, post.ID, post.Title, post.Text, post.UserID, post.UserName, post.UserAvatar, post.ImageURL, post.CreatedAt, post.UpdatedAt).Scan(&post.ID)
	if err != nil {
		return nil, fmt.Errorf("error creating post: %v", err)
	}

	return post, nil
}

// Обновление поста
func (r *PostRepositoryPg) UpdatePost(post *models.Post) (*models.Post, error) {
	// Запрос на обновление поста
	query := `UPDATE posts SET title = $1, text = $2, user_name = $3, user_avatar = $4, image_url = $5, updated_at = $6 WHERE id = $7`

	_, err := r.db.Exec(query, post.Title, post.Text, post.UserName, post.UserAvatar, post.ImageURL, post.UpdatedAt, post.ID)
	if err != nil {
		return nil, fmt.Errorf("error updating post: %v", err)
	}

	return post, nil
}

// Удаление поста
func (r *PostRepositoryPg) DeletePost(id string) error {
	// Запрос на удаление поста
	query := `DELETE FROM posts WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting post: %v", err)
	}

	return nil
}

// Получение поста по ID
func (r *PostRepositoryPg) GetPostByID(id string) (*models.Post, error) {
	// Запрос на получение поста
	query := `SELECT id, title, text, user_id, user_name, user_avatar, image_url, created_at, updated_at FROM posts WHERE id = $1`

	var post models.Post
	err := r.db.QueryRow(query, id).Scan(&post.ID, &post.Title, &post.Text, &post.UserID, &post.UserName, &post.UserAvatar, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting post by id: %v", err)
	}

	return &post, nil
}

// Получение всех постов
func (r *PostRepositoryPg) GetAllPosts() ([]*models.Post, error) {
	// Запрос на получение всех постов
	query := `SELECT id, title, text, user_id, user_name, user_avatar, image_url, created_at, updated_at FROM posts`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting all posts: %v", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Text, &post.UserID, &post.UserName, &post.UserAvatar, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		posts = append(posts, &post)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %v", rows.Err())
	}

	return posts, nil
}
