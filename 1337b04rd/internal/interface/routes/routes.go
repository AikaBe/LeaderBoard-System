package routes

import (
	"net/http"

	"1337b04rd/internal/interface/handlers"
	"1337b04rd/internal/interface/middleware"
)

func RegisterRoutes(mux *http.ServeMux, postHandler *handlers.PostHandler, commentHandler *handlers.CommentHandler, authMiddleware middleware.AuthMiddleware) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/templates"))))

	mux.Handle("/posts", authMiddleware.LoginOrLastVisitHandler(http.HandlerFunc(postHandler.GetAllPosts)))
	mux.Handle("/create-post", authMiddleware.LoginOrLastVisitHandler(http.HandlerFunc(postHandler.SubmitPost)))
	mux.Handle("/create", authMiddleware.LoginOrLastVisitHandler(http.HandlerFunc(postHandler.ServeCreatePostForm)))
	mux.Handle("/post/", authMiddleware.LoginOrLastVisitHandler(http.HandlerFunc(postHandler.GetPostByID)))

	mux.Handle("/submit-comment", authMiddleware.LoginOrLastVisitHandler(http.HandlerFunc(commentHandler.CreateComment)))

	mux.Handle("/archive", authMiddleware.LoginOrLastVisitHandler(http.HandlerFunc(postHandler.GetArchivedPostsHandler)))
	mux.Handle("/archived/post/", authMiddleware.LoginOrLastVisitHandler(http.HandlerFunc(postHandler.GetArchivedPostByID)))
}
