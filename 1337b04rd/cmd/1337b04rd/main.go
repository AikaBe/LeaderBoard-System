package main

import (
	"1337b04rd/internal/adapters/database"
	d "1337b04rd/internal/adapters/database"
	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/app/domain/ports"
	"1337b04rd/internal/app/domain/services"
	"1337b04rd/internal/interface/handlers"
	"1337b04rd/internal/interface/middleware"
	"1337b04rd/internal/interface/routes"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func initRepository(dsn string) (ports.PostRepository, ports.CommentRepository, *sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, nil, err
	}

	postRepo := d.NewPostRepositoryPg(db)
	commentRepo := d.NewCommentRepositoryPg(db)

	return postRepo, commentRepo, db, nil
}

func main() {
	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dsn := "host=db port=5432 user=board_user password=board_pass dbname=board_db sslmode=disable"

	// Connect to DB
	postRepo, commentRepo, db, err := initRepository(dsn)
	if err != nil {
		logger.Error("Failed to initialize repository", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create services
	commentService := services.NewCommentService(commentRepo)
	sessionRepo := &database.PostgresSessionRepo{DB: db}
	sessionService := services.NewSessionService(sessionRepo)
	postService := services.NewPostService(postRepo, sessionRepo)
	postService.StartArchiver()

	// S3 Adapter setup
	s3Adapter := &s3.Adapter{
		TripleSBaseURL:  "http://triple-s:9000",
		PublicAccessURL: "http://localhost:9000",
	}

	// Handlers
	postHandler := handlers.NewPostHandler(postService, s3Adapter, commentService)
	commentHandler := handlers.NewCommentHandler(commentService, sessionRepo, s3Adapter)
	authMiddleware := middleware.AuthMiddleware{SessionService: sessionService}

	// Router setup
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, postHandler, commentHandler, authMiddleware)

	// Start server
	logger.Info("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("Failed to start server", "error", err)
	}
}
