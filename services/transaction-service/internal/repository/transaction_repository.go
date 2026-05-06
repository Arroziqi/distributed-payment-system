package repository

import (
	"context"
	"time"

	"transaction-service/internal/domain"
)

type TransactionRepository interface {
	ProcessPayment(ctx context.Context, tx domain.Transaction, idemKey, requestHash string, expiresAt time.Time) (domain.Transaction, []byte, bool, error)
	List(ctx context.Context, limit, offset int) ([]domain.Transaction, error)
	GetByID(ctx context.Context, id string) (domain.Transaction, error)
}
