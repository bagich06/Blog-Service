package service

import (
	"blog/internal/models"
	"blog/internal/repository"
	"context"
	"fmt"
	"time"
)

type LoginLimiterService struct {
	redisRepo *repository.RedisLoginLimiter
	limits    map[string]*models.LoginLimit
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

	attempts, err := s.redisRepo.GetLoginAttempts(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get login attempts: %w", err)
	}

	limit := s.limits["login"]
	if attempts >= limit.MaxAttempts {
		err = s.redisRepo.BlockUser(ctx, userID, limit.Window)
		if err != nil {
			return fmt.Errorf("failed to block user: %w", err)
		}
		return fmt.Errorf("too many login attempts, user blocked for %v", limit.Window)
	}

	return nil
}

func (s *LoginLimiterService) RecordFailedLogin(ctx context.Context, userID int) error {
	limit := s.limits["login"]
	_, err := s.redisRepo.IncrementLoginAttempts(ctx, userID, limit.Window)
	return err
}

func (s *LoginLimiterService) RecordSuccessfulLogin(ctx context.Context, userID int) error {
	return s.redisRepo.ResetLoginAttempts(ctx, userID)
}
