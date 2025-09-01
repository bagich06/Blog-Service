package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCommentsLimiter struct {
	client *redis.Client
}

func NewRedisCommentsLimiter(addrStr string) (*RedisCommentsLimiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addrStr,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCommentsLimiter{client: client}, nil
}

func (l *RedisCommentsLimiter) GetCommentAttempts(ctx context.Context, userID int, postID int) (int, error) {
	key := fmt.Sprintf("user:%d:post:%d:comment_attempts", userID, postID)
	res, err := l.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}

	var attempts int
	_, err = fmt.Sscanf(res, "%d", &attempts)
	return attempts, err
}

func (l *RedisCommentsLimiter) IncrementCommentAttempts(ctx context.Context, userID int, postID int, window time.Duration) (int, error) {
	key := fmt.Sprintf("user:%d:post:%d:comment_attempts", userID, postID)

	attempts, err := l.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if attempts == 1 {
		l.client.Expire(ctx, key, window)
	}

	return int(attempts), nil
}

func (l *RedisCommentsLimiter) BlockUserComment(ctx context.Context, userID int, postID int, duration time.Duration) error {
	key := fmt.Sprintf("user:%d:post:%d:blocked", userID, postID)
	return l.client.Set(ctx, key, "1", duration).Err()
}

func (l *RedisCommentsLimiter) IsUserBlockedComment(ctx context.Context, userID int, postID int) (bool, error) {
	key := fmt.Sprintf("user:%d:post:%d:blocked", userID, postID)
	exists, err := l.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (l *RedisCommentsLimiter) ResetCommentAttempts(ctx context.Context, userID int, postID int) error {
	key := fmt.Sprintf("user:%d:post:%d:comment_attempts", userID, postID)
	return l.client.Del(ctx, key).Err()
}
