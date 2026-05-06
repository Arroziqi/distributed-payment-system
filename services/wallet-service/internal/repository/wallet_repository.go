package repository

import (
	"context"

	"wallet-service/internal/domain"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet domain.Wallet) error
	GetByUserID(ctx context.Context, userID string) (domain.Wallet, error)
	Topup(ctx context.Context, userID string, amount int64, txnID string) (domain.Wallet, error)
	Withdraw(ctx context.Context, userID string, amount int64, txnID string) (domain.Wallet, error)
	Transfer(ctx context.Context, fromUserID string, toUserID string, amount int64, txnID string) (domain.Wallet, domain.Wallet, error)
}
