package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"

	"1337b04rd/internal/app/domain/models"
	"1337b04rd/internal/app/domain/ports"
)

// CommentRepositoryPg represents a repository for managing comments in PostgreSQL.
type CommentRepositoryPg struct {
	db *sql.DB
}

// NewCommentRepositoryPg creates a new instance of the comment repository.
func NewCommentRepositoryPg(db *sql.DB) ports.CommentRepository {
	return &CommentRepositoryPg{db: db}
}

// CreateComment creates a new comment for the given post.
func (r *CommentRepositoryPg) CreateComment(comment models.Comment) (*models.Comment, error) {
	// Check if the post exists
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`
	err := r.db.QueryRow(query, comment.PostID).Scan(&exists)
	if err != nil {
		slog.Error("Error checking post existence", "postID", comment.PostID, "error", err)
		return nil, fmt.Errorf("error checking post existence: %v", err)
	}

	if !exists {
		slog.Warn("Post does not exist", "postID", comment.PostID)
		return nil, fmt.Errorf("post with id %d does not exist", comment.PostID)
	}

	// Insert the new comment
	query = `
	INSERT INTO comments (post_id, parent_comment_id, user_name, user_avatar, text, image_url, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`

	var id int
	err = r.db.QueryRow(
		query,
		comment.PostID,
		comment.ParentCommentID,
		comment.UserName,
		comment.UserAvatar,
		comment.Text,
		comment.ImageURL, // Added field for image
		comment.CreatedAt,
	).Scan(&id)
	if err != nil {
		slog.Error("Error creating comment", "error", err)
		return nil, fmt.Errorf("error creating comment: %v", err)
	}

	comment.ID = id
	slog.Info("Successfully created comment", "commentID", id)
	return &comment, nil
}

// GetCommentsByPostID retrieves all comments for the given post.
func (r *CommentRepositoryPg) GetCommentsByPostID(postID int) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, parent_comment_id, user_name, user_avatar, text, image_url, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		slog.Error("Error getting comments", "postID", postID, "error", err)
		return nil, fmt.Errorf("error getting comments: %v", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	commentMap := make(map[int][]*models.Comment)

	// Reading the comments
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(
			&c.ID,
			&c.PostID,
			&c.ParentCommentID,
			&c.UserName,
			&c.UserAvatar,
			&c.Text,
			&c.ImageURL, // Added field for image
			&c.CreatedAt,
		)
		if err != nil {
			slog.Error("Error scanning comment", "error", err)
			return nil, fmt.Errorf("error scanning comment: %v", err)
		}

		// If it's a reply, add it to the map
		if c.ParentCommentID != nil {
			commentMap[*c.ParentCommentID] = append(commentMap[*c.ParentCommentID], &c)
		} else {
			comments = append(comments, &c)
		}
	}

	// Adding replies to parent comments
	for i := range comments {
		parent := comments[i]
		if replies, ok := commentMap[parent.ID]; ok {
			parent.Replies = replies
		}
	}

	slog.Info("Successfully retrieved comments", "postID", postID, "commentCount", len(comments))
	return comments, nil
}

// DeleteComment deletes a comment by its ID.
func (r *CommentRepositoryPg) DeleteComment(commentID string) error {
	// Convert commentID from string to int
	idInt, err := strconv.Atoi(commentID)
	if err != nil {
		slog.Error("Invalid comment ID", "commentID", commentID, "error", err)
		return fmt.Errorf("invalid comment id: %v", err)
	}

	query := `DELETE FROM comments WHERE id = $1`
	_, err = r.db.Exec(query, idInt)
	if err != nil {
		slog.Error("Error deleting comment", "commentID", idInt, "error", err)
		return fmt.Errorf("error deleting comment: %v", err)
	}

	slog.Info("Successfully deleted comment", "commentID", idInt)
	return nil
}
