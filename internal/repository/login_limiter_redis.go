package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisLoginLimiter struct {
	client *redis.Client
}

func NewRedisLoginLimiter(addrStr string) (*RedisLoginLimiter, error) {
	clinet := redis.NewClient(&redis.Options{
		Addr: addrStr,
	})

	ctx := context.Background()
	_, err := clinet.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisLoginLimiter{client: clinet}, nil
}

func (l *RedisLoginLimiter) GetLoginAttempts(ctx context.Context, userID int) (int, error) {
	key := fmt.Sprintf("user:%d:login_attempts", userID)
	res, err := l.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}

	var attempts int
	_, err = fmt.Sscanf(res, "%d", &attempts)
	return attempts, err
}

func (l *RedisLoginLimiter) IncrementLoginAttempts(ctx context.Context, userID int, window time.Duration) (int, error) {
	key := fmt.Sprintf("user:%d:login_attempts", userID)

	attempts, err := l.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if attempts == 1 {
		l.client.Expire(ctx, key, window)
	}

	return int(attempts), nil
}

func (l *RedisLoginLimiter) BlockUser(ctx context.Context, userID int, duration time.Duration) error {
	key := fmt.Sprintf("user:%d:blocked", userID)
	return l.client.Set(ctx, key, "1", duration).Err()
}

func (l *RedisLoginLimiter) IsUserBlocked(ctx context.Context, userID int) (bool, error) {
	key := fmt.Sprintf("user:%d:blocked", userID)
	exists, err := l.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (l *RedisLoginLimiter) ResetLoginAttempts(ctx context.Context, userID int) error {
	key := fmt.Sprintf("user:%d:login_attempts", userID)
	return l.client.Del(ctx, key).Err()
}
