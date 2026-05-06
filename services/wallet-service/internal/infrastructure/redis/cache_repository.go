package redis

import (
	"context"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type BalanceCacheRepository struct {
	client *goredis.Client
}

func NewBalanceCacheRepository(client *goredis.Client) BalanceCacheRepository {
	return BalanceCacheRepository{client: client}
}

func (r BalanceCacheRepository) Get(ctx context.Context, userID string) (int64, bool, error) {
	v, err := r.client.Get(ctx, key(userID)).Result()
	if err != nil {
		if err == goredis.Nil {
			return 0, false, nil
		}
		return 0, false, err
	}
	parsed, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false, err
	}
	return parsed, true, nil
}

func (r BalanceCacheRepository) Set(ctx context.Context, userID string, balance int64, ttl time.Duration) error {
	return r.client.Set(ctx, key(userID), strconv.FormatInt(balance, 10), ttl).Err()
}

func (r BalanceCacheRepository) Invalidate(ctx context.Context, userID string) error {
	return r.client.Del(ctx, key(userID)).Err()
}

func key(userID string) string {
	return "wallet:balance:" + userID
}
