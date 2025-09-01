package service

import (
	"blog/internal/models"
	"blog/internal/repository"
	"context"
	"fmt"
	"time"
)

type LoginLimiterService struct {
	redisRepo *repository.RedisLoginLimiter // слой для работы с Redis
	limits    map[string]*models.LoginLimit // конфигурация лимитов
}

func NewLoginLimiterService(redisRepo *repository.RedisLoginLimiter) *LoginLimiterService {
	return &LoginLimiterService{
		redisRepo: redisRepo,
		limits: map[string]*models.LoginLimit{
			"login": {
				MaxAttempts: 3,
				Window:      1 * time.Hour,
			},
		},
		// конфигурация лимитов по ключу "login"
	}
}

func (s *LoginLimiterService) CheckLoginLimit(ctx context.Context, userID int) error {
	blocked, err := s.redisRepo.IsUserBlocked(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to check if user blocked: %w", err)
	}

	if blocked {
		return fmt.Errorf("user is blocked due to too many failed login attempts")
	}

	// Проверяем заблокирован ли пользователь

	attempts, err := s.redisRepo.GetLoginAttempts(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get login attempts: %w", err)
	}

	// Получаем кол-во попыток для входа

	limit := s.limits["login"] // в переменную limit передаем конфигурацию login MaxAttempts and Window
	if attempts >= limit.MaxAttempts {
		err = s.redisRepo.BlockUser(ctx, userID, limit.Window) // блокируем, если попытки превышают MaxAttempts
		if err != nil {
			return fmt.Errorf("failed to block user: %w", err)
		}
		return fmt.Errorf("too many login attempts, user blocked for %v", limit.Window)
	}

	return nil
}

func (s *LoginLimiterService) RecordFailedLogin(ctx context.Context, userID int) error {
	limit := s.limits["login"]
	// в переменную limit передаем конфигурацию login MaxAttempts and Window
	_, err := s.redisRepo.IncrementLoginAttempts(ctx, userID, limit.Window)
	// инкрементим переменную
	return err
}

func (s *LoginLimiterService) RecordSuccessfulLogin(ctx context.Context, userID int) error {
	return s.redisRepo.ResetLoginAttempts(ctx, userID)
	// сбрасываем попытки на вход
}
