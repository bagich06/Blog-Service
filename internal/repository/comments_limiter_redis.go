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

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCommentsLimiter{client: client}, nil
}

func (l *RedisCommentsLimiter) GetCommentAttempts(ctx context.Context, userID int, postID int) (int, error) {
	key := fmt.Sprintf("user:%d:post:%d:comment_attempts", userID, postID) // устанавливаем ключ с id поста и id юзера
	res, err := l.client.Get(ctx, key).Result()                            // получаем значение из Redis
	if err == redis.Nil {
		return 0, nil // если нет ключа значит ретерним 0
	}

	var attempts int
	_, err = fmt.Sscanf(res, "%d", &attempts)
	return attempts, err
	// записываем в переменную attmepts и возвращаем
}

func (l *RedisCommentsLimiter) IncrementCommentAttempts(ctx context.Context, userID int, postID int, window time.Duration) (int, error) {
	key := fmt.Sprintf("user:%d:post:%d:comment_attempts", userID, postID) // задаем ключ

	attempts, err := l.client.Incr(ctx, key).Result() // инкрементим значение, если ключ не найден устанавливаем 0
	if err != nil {
		return 0, err
	}

	if attempts == 1 {
		l.client.Expire(ctx, key, window) // когда была совершена 1ая попытка начинаем отсчет промежутка
	}

	return int(attempts), nil
	// возвращаем попытки
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
