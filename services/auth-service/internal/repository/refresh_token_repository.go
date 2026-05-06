package repository

import (
	"context"
	"time"
)

type RefreshTokenRepository interface {
	Store(ctx context.Context, token string, userID string, ttl time.Duration) error
	Consume(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
}
