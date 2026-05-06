package repository

import (
	"context"
	"time"
)

type BalanceCacheRepository interface {
	Get(ctx context.Context, userID string) (int64, bool, error)
	Set(ctx context.Context, userID string, balance int64, ttl time.Duration) error
	Invalidate(ctx context.Context, userID string) error
}
