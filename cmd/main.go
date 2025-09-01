package main

import (
	"blog/internal/api"
	"blog/internal/repository"
	"blog/internal/service"
	"log"

	"github.com/gorilla/mux"
)

func main() {
	// Подключение к PostgreSQL
	db, err := repository.New("postgres://postgres:postgres@localhost:5432/blog_db")
	if err != nil {
		log.Fatal(err)
	}

	// Создание Redis клиентов для rate limiting
	loginLimiter, err := repository.NewRedisLoginLimiter("localhost:6379")
	if err != nil {
		log.Fatal("Failed to create login limiter:", err)
	}

	commentLimiter, err := repository.NewRedisCommentsLimiter("localhost:6379")
	if err != nil {
		log.Fatal("Failed to create comment limiter:", err)
	}

	// Создание сервисов
	loginLimiterService := service.NewLoginLimiterService(loginLimiter)
	commentLimiterService := service.NewCommentLimiterService(commentLimiter)

	// Создание API с rate limiting сервисами
	api := api.New(mux.NewRouter(), db, loginLimiterService, commentLimiterService)
	api.Handle()
	log.Fatal(api.ListenAndServe("localhost:8080"))
}
