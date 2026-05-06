package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"auth-service/internal/repository"

	goredis "github.com/redis/go-redis/v9"
)

type RefreshTokenRepository struct {
	client    *goredis.Client
	keyPrefix string
}

func NewRefreshTokenRepository(client *goredis.Client, keyPrefix string) RefreshTokenRepository {
	return RefreshTokenRepository{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

func (r RefreshTokenRepository) Store(ctx context.Context, token string, userID string, ttl time.Duration) error {
	return r.client.Set(ctx, r.key(token), userID, ttl).Err()
}

func (r RefreshTokenRepository) Consume(ctx context.Context, token string) (string, error) {
	result, err := r.client.GetDel(ctx, r.key(token)).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return "", repository.ErrNotFound
		}
		return "", fmt.Errorf("redis getdel: %w", err)
	}
	return result, nil
}

func (r RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	return r.client.Del(ctx, r.key(token)).Err()
}

func (r RefreshTokenRepository) key(token string) string {
	return r.keyPrefix + token
}
