package handlers

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"1337b04rd/internal/app/domain/ports"
	"1337b04rd/internal/app/domain/services"
)

type PostHandler struct {
	PostService    *services.PostService
	S3Adapter      ports.S3Adapter
	CommentService ports.CommentService
}

func NewPostHandler(postService *services.PostService, s3Adapter ports.S3Adapter, commentService ports.CommentService) *PostHandler {
	return &PostHandler{
		PostService:    postService,
		S3Adapter:      s3Adapter,
		CommentService: commentService,
	}
}

func (h *PostHandler) ServeCreatePostForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/create-post.html"))
	tmpl.Execute(w, nil)
}

func (h *PostHandler) SubmitPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		slog.Error("Failed to parse form", "error", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	sessionIDRaw := r.Context().Value("sessionId")
	sessionID, ok := sessionIDRaw.(string)
	if !ok || sessionID == "" {
		slog.Warn("Unauthorized access: missing or invalid sessionId")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("subject")
	text := r.FormValue("comment")

	imageURL, err := h.S3Adapter.UploadImage(r, "post")
	if err != nil {
		slog.Error("Image upload failed", "error", err)
		http.Error(w, "Image upload failed", http.StatusInternalServerError)
		return
	}

	createdPost, err := h.PostService.CreatePost(sessionID, title, text, imageURL)
	if err != nil {
		slog.Error("Failed to create post", "error", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	slog.Info("Post created", "post_id", createdPost.ID)
	http.Redirect(w, r, "/create", http.StatusSeeOther)
}

func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostService.GetAllPosts()
	if err != nil {
		slog.Error("Failed to fetch posts", "error", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/catalog.html"))
	if err := tmpl.Execute(w, posts); err != nil {
		slog.Error("Failed to render catalog template", "error", err)
		http.Error(w, "Failed to render posts", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/post/")
	if id == "" {
		http.Error(w, "Missing post ID", http.StatusBadRequest)
		return
	}

	post, err := h.PostService.GetPostByID(id)
	if err != nil {
		slog.Error("Failed to fetch post by ID", "post_id", id, "error", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.NotFound(w, r)
		return
	}

	comments, err := h.CommentService.GetCommentsByPostID(post.ID)
	if err != nil {
		slog.Error("Failed to load comments", "post_id", post.ID, "error", err)
		http.Error(w, "Error loading comments", http.StatusInternalServerError)
		return
	}
	post.Comments = comments

	tmpl := template.Must(template.ParseFiles("web/templates/post.html"))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, post); err != nil {
		slog.Error("Failed to render post template", "post_id", post.ID, "error", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func (h *PostHandler) GetArchivedPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostService.GetArchivedPosts()
	if err != nil {
		slog.Error("Failed to fetch archived posts", "error", err)
		http.Error(w, "Failed to fetch archived posts", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/archive.html"))
	if err := tmpl.Execute(w, posts); err != nil {
		slog.Error("Failed to render archive template", "error", err)
		http.Error(w, "Failed to render posts", http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetArchivedPostByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/archived/post/")
	if id == "" {
		http.Error(w, "Missing post ID", http.StatusBadRequest)
		return
	}

	post, err := h.PostService.GetArchivedPostByID(id)
	if err != nil {
		slog.Error("Failed to fetch archived post", "post_id", id, "error", err)
		http.Error(w, "Failed to fetch archived post", http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.NotFound(w, r)
		return
	}

	comments, err := h.CommentService.GetCommentsByPostID(post.ID)
	if err != nil {
		slog.Error("Failed to load archived comments", "post_id", post.ID, "error", err)
		http.Error(w, "Error loading comments", http.StatusInternalServerError)
		return
	}
	post.Comments = comments

	tmpl := template.Must(template.ParseFiles("web/templates/archive-post.html"))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, post); err != nil {
		slog.Error("Failed to render archived post template", "post_id", post.ID, "error", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
