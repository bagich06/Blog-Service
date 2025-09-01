package service

import (
	"blog/internal/models"
	"blog/internal/repository"
	"context"
	"fmt"
	"time"
)

type CommentLimiterService struct {
	redisRepo *repository.RedisCommentsLimiter
	limits    map[string]*models.CommentLimit
}

func NewCommentLimiterService(redisRepo *repository.RedisCommentsLimiter) *CommentLimiterService {
	return &CommentLimiterService{
		redisRepo: redisRepo,
		limits: map[string]*models.CommentLimit{
			"comment": {
				MaxAttempts: 5,
				Window:      1 * time.Hour,
			},
		},
	}
}

func (s *CommentLimiterService) CheckCommentLimit(ctx context.Context, userID int, postID int) error {
	blocked, err := s.redisRepo.IsUserBlockedComment(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to check if user blocked: %w", err)
	}

	if blocked {
		return fmt.Errorf("user is blocked from commenting on this post due to too many attempts")
	}

	attempts, err := s.redisRepo.GetCommentAttempts(ctx, userID, postID)
	if err != nil {
		return fmt.Errorf("failed to get comment attempts: %w", err)
	}

	limit := s.limits["comment"]
	if attempts >= limit.MaxAttempts {
		err = s.redisRepo.BlockUserComment(ctx, userID, postID, limit.Window)
		if err != nil {
			return fmt.Errorf("failed to block user: %w", err)
		}
		return fmt.Errorf("too many comment attempts on this post, user blocked for %v", limit.Window)
	}

	return nil
}

func (s *CommentLimiterService) RecordCommentAttempt(ctx context.Context, userID int, postID int) error {
	limit := s.limits["comment"]
	_, err := s.redisRepo.IncrementCommentAttempts(ctx, userID, postID, limit.Window)
	return err
}

func (s *CommentLimiterService) ResetCommentAttempts(ctx context.Context, userID int, postID int) error {
	return s.redisRepo.ResetCommentAttempts(ctx, userID, postID)
}
