package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
)

type PostRepositoryPg struct {
	db *sql.DB
}

func NewPostRepositoryPg(db *sql.DB) ports.PostRepository {
	return &PostRepositoryPg{db: db}
}

// CreatePost creates a new post and returns the created post with its ID.
func (r *PostRepositoryPg) CreatePost(post *models.Post) (*models.Post, error) {
	query := `INSERT INTO posts (title, text, user_name, user_avatar, image_url, created_at, updated_at, is_hidden) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	// Execute the query and get the automatically generated ID
	err := r.db.QueryRow(query, post.Title, post.Text, post.UserName, post.UserAvatar, post.ImageURL, post.CreatedAt, post.UpdatedAt, post.IsHidden).Scan(&post.ID)
	if err != nil {
		slog.Error("Error creating post", "error", err)
		return nil, fmt.Errorf("error creating post: %v", err)
	}

	slog.Info("Successfully created post", "postID", post.ID)
	return post, nil
}

// GetAllPosts retrieves all posts that are not hidden and archived.
func (r *PostRepositoryPg) GetAllPosts() ([]*models.Post, error) {
	query := `SELECT id, title, text,  user_name, user_avatar, image_url, created_at, updated_at, is_hidden 
	          FROM posts 
	          WHERE is_hidden = FALSE AND archived_at IS NULL`

	rows, err := r.db.Query(query)
	if err != nil {
		slog.Error("Error getting all posts", "error", err)
		return nil, fmt.Errorf("error getting all posts: %v", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Text, &post.UserName, &post.UserAvatar, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt, &post.IsHidden)
		if err != nil {
			slog.Error("Error scanning row", "error", err)
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		posts = append(posts, &post)
	}

	if rows.Err() != nil {
		slog.Error("Error iterating rows", "error", rows.Err())
		return nil, fmt.Errorf("error iterating rows: %v", rows.Err())
	}

	slog.Info("Successfully retrieved posts", "postCount", len(posts))
	return posts, nil
}

// GetPostByID retrieves a post by its ID.
func (r *PostRepositoryPg) GetPostByID(id string) (*models.Post, error) {
	// Convert ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("Invalid ID format", "id", id, "error", err)
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	query := `SELECT id, title, text,  user_name, user_avatar, image_url, created_at, updated_at, is_hidden 
	          FROM posts WHERE id = $1 AND archived_at IS NULL`

	var post models.Post
	err = r.db.QueryRow(query, idInt).Scan(&post.ID, &post.Title, &post.Text, &post.UserName, &post.UserAvatar, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt, &post.IsHidden)
	if err != nil {
		slog.Error("Error getting post by ID", "id", idInt, "error", err)
		return nil, fmt.Errorf("error getting post by id: %v", err)
	}

	slog.Info("Successfully retrieved post by ID", "postID", post.ID)
	return &post, nil
}

// ArchivePost archives a post by setting its archived_at field.
func (r *PostRepositoryPg) ArchivePost(postID int) error {
	query := `UPDATE posts SET archived_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(context.Background(), query, time.Now(), postID)
	if err != nil {
		slog.Error("Error archiving post", "postID", postID, "error", err)
		return err
	}

	slog.Info("Successfully archived post", "postID", postID)
	return nil
}

// GetArchivedPosts retrieves all archived posts.
func (r *PostRepositoryPg) GetArchivedPosts() ([]*models.Post, error) {
	query := `SELECT id, title, text, user_name, user_avatar, image_url, created_at, updated_at, archived_at
	          FROM posts WHERE archived_at IS NOT NULL ORDER BY archived_at DESC`

	rows, err := r.db.QueryContext(context.Background(), query)
	if err != nil {
		slog.Error("Error getting archived posts", "error", err)
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Text, &post.UserName, &post.UserAvatar, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt, &post.ArchivedAt)
		if err != nil {
			slog.Error("Error scanning archived post", "error", err)
			return nil, err
		}
		posts = append(posts, &post)
	}

	slog.Info("Successfully retrieved archived posts", "postCount", len(posts))
	return posts, nil
}

// GetArchivedPostByID retrieves an archived post by its ID.
func (r *PostRepositoryPg) GetArchivedPostByID(id string) (*models.Post, error) {
	// Convert ID from string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		slog.Error("Invalid ID format", "id", id, "error", err)
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	query := `SELECT id, title, text, user_name, user_avatar, image_url, created_at, updated_at, is_hidden
	          FROM posts WHERE id = $1 AND archived_at IS NOT NULL`

	var post models.Post
	err = r.db.QueryRow(query, idInt).Scan(&post.ID, &post.Title, &post.Text, &post.UserName, &post.UserAvatar, &post.ImageURL, &post.CreatedAt, &post.UpdatedAt, &post.IsHidden)
	if err != nil {
		slog.Error("Error getting archived post by ID", "id", idInt, "error", err)
		return nil, fmt.Errorf("error getting archived post by id: %v", err)
	}

	slog.Info("Successfully retrieved archived post by ID", "postID", post.ID)
	return &post, nil
}
