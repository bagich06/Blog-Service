package api

import (
	"blog/internal/middleware"
	"blog/internal/repository"
	"blog/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

type api struct {
	r              *mux.Router
	db             *repository.PGRepo
	loginLimiter   *service.LoginLimiterService
	commentLimiter *service.CommentLimiterService
}

func New(r *mux.Router, db *repository.PGRepo, loginLimiter *service.LoginLimiterService, commentLimiter *service.CommentLimiterService) *api {
	return &api{
		r:              r,
		db:             db,
		loginLimiter:   loginLimiter,
		commentLimiter: commentLimiter,
	}
}

func (api *api) Handle() {
	api.r.HandleFunc("/api/login", api.LoginHandler)
	api.r.HandleFunc("/api/register", api.RegisterHandler)
	api.r.HandleFunc("/api/user/create", api.CreateUser)
	api.r.HandleFunc("/api/user/get", api.GetUserByID).Queries("id", "{id}")
	api.r.HandleFunc("/api/user/delete", api.DeleteUser).Queries("id", "{id}")
	api.r.HandleFunc("/api/user/get", api.GetUsers)
	api.r.HandleFunc("/api/post/create", middleware.AuthMiddleware(api.CreatePost))
	api.r.HandleFunc("/api/post/get", middleware.AuthMiddleware(api.GetPostById)).Queries("id", "{id}")
	api.r.HandleFunc("/api/post/get", middleware.AuthMiddleware(api.GetAllPosts))
	api.r.HandleFunc("/api/post/delete", middleware.AuthMiddleware(api.DeletePostById)).Queries("id", "{id}")
	api.r.HandleFunc("/api/comment/create", middleware.AuthMiddleware(api.CreateComment)).Queries("post_id", "{post_id}")
}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api.r)
}
