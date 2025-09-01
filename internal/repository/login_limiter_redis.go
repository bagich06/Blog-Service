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
	client := redis.NewClient(&redis.Options{
		Addr: addrStr,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisLoginLimiter{client: client}, nil
}

func (l *RedisLoginLimiter) GetLoginAttempts(ctx context.Context, userID int) (int, error) {
	key := fmt.Sprintf("user:%d:login_attempts", userID) // формирование ключа для каждого пользователя
	res, err := l.client.Get(ctx, key).Result()          // получаем значение из Redis
	if err == redis.Nil {
		return 0, nil
	}
	// обработка случая, когда ключ не найден

	var attempts int
	_, err = fmt.Sscanf(res, "%d", &attempts) // строку в число
	return attempts, err
}

func (l *RedisLoginLimiter) IncrementLoginAttempts(ctx context.Context, userID int, window time.Duration) (int, error) {
	key := fmt.Sprintf("user:%d:login_attempts", userID) // формирование ключа для каждого пользователя

	attempts, err := l.client.Incr(ctx, key).Result() // читаем значение ключа и инкрементим, если не найден то 0
	if err != nil {
		return 0, err
	}

	if attempts == 1 {
		l.client.Expire(ctx, key, window)
	}
	// когда получаем 1 попытку то мы устанавливаем время чтобы указат промежуток в который происходит попытка входа

	return int(attempts), nil
}

func (l *RedisLoginLimiter) BlockUser(ctx context.Context, userID int, duration time.Duration) error {
	key := fmt.Sprintf("user:%d:blocked", userID)      // создаем новый ключ
	return l.client.Set(ctx, key, "1", duration).Err() // устанавливаем значение 1 для блокировки
}

func (l *RedisLoginLimiter) IsUserBlocked(ctx context.Context, userID int) (bool, error) {
	key := fmt.Sprintf("user:%d:blocked", userID)
	exists, err := l.client.Exists(ctx, key).Result() // проверяем наличие блока по ключу (true/false)
	if err != nil {
		return false, err
	}
	return exists > 0, nil // exists > 0 - если exists = 1 → true, если exists = 0 → false
}

func (l *RedisLoginLimiter) ResetLoginAttempts(ctx context.Context, userID int) error {
	key := fmt.Sprintf("user:%d:login_attempts", userID) // устанавливаем ключ
	return l.client.Del(ctx, key).Err()                  // Удаляем его и получаем ошибку
}
