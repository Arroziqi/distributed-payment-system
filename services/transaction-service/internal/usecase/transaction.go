package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"transaction-service/internal/domain"
	"transaction-service/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrIdempotencyRequired = errors.New("idempotency key is required")
	ErrIdempotencyMismatch = errors.New("idempotency key mismatch")
	ErrIdempotencyBusy     = errors.New("idempotency key is being processed")
	ErrTransactionNotFound = errors.New("transaction not found")
)

type TransactionUsecase struct {
	repo    repository.TransactionRepository
	pub     repository.EventPublisher
	idemTTL time.Duration
}

func NewTransactionUsecase(repo repository.TransactionRepository, pub repository.EventPublisher, idemTTL time.Duration) *TransactionUsecase {
	return &TransactionUsecase{repo: repo, pub: pub, idemTTL: idemTTL}
}

type ProcessPaymentInput struct {
	IdempotencyKey string `json:"-"`
	ExternalID     string `json:"external_id"`
	Type           string `json:"type"`
	FromWalletID   string `json:"from_wallet_id"`
	ToWalletID     string `json:"to_wallet_id"`
	Amount         int64  `json:"amount"`
}

func (u *TransactionUsecase) ProcessPayment(ctx context.Context, in ProcessPaymentInput) (domain.Transaction, bool, error) {
	if in.IdempotencyKey == "" {
		return domain.Transaction{}, false, ErrIdempotencyRequired
	}
	if in.Type == "" || in.Amount <= 0 || in.FromWalletID == "" || in.ToWalletID == "" || in.FromWalletID == in.ToWalletID {
		return domain.Transaction{}, false, ErrInvalidInput
	}

	reqHash, err := hashInput(in)
	if err != nil {
		return domain.Transaction{}, false, fmt.Errorf("hash input: %w", err)
	}

	now := time.Now().UTC()
	tx := domain.Transaction{
		ID:           uuid.NewString(),
		ExternalID:   in.ExternalID,
		Type:         in.Type,
		FromWalletID: in.FromWalletID,
		ToWalletID:   in.ToWalletID,
		Amount:       in.Amount,
		Status:       "completed",
		Version:      0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	storedTx, snapshot, replayed, err := u.repo.ProcessPayment(ctx, tx, in.IdempotencyKey, reqHash, now.Add(u.idemTTL))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrIdempotencyMismatch):
			return domain.Transaction{}, false, ErrIdempotencyMismatch
		case errors.Is(err, repository.ErrIdempotencyBusy):
			return domain.Transaction{}, false, ErrIdempotencyBusy
		default:
			return domain.Transaction{}, false, err
		}
	}

	if replayed && len(snapshot) > 0 {
		var prev domain.Transaction
		if err := json.Unmarshal(snapshot, &prev); err == nil {
			return prev, true, nil
		}
	}

	if err := u.publishCompleted(ctx, storedTx); err != nil {
		return domain.Transaction{}, false, fmt.Errorf("publish transaction completed: %w", err)
	}
	return storedTx, replayed, nil
}

func (u *TransactionUsecase) ListHistory(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return u.repo.List(ctx, limit, offset)
}

func (u *TransactionUsecase) GetByID(ctx context.Context, id string) (domain.Transaction, error) {
	if id == "" {
		return domain.Transaction{}, ErrInvalidInput
	}
	tx, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Transaction{}, ErrTransactionNotFound
		}
		return domain.Transaction{}, err
	}
	return tx, nil
}

func (u *TransactionUsecase) publishCompleted(ctx context.Context, tx domain.Transaction) error {
	body, err := json.Marshal(map[string]any{
		"transaction_id": tx.ID,
		"type":           tx.Type,
		"status":         tx.Status,
		"amount":         tx.Amount,
		"from_wallet_id": tx.FromWalletID,
		"to_wallet_id":   tx.ToWalletID,
		"occurred_at":    tx.UpdatedAt,
	})
	if err != nil {
		return err
	}
	return u.pub.Publish(ctx, "transaction.completed", body)
}

func hashInput(in ProcessPaymentInput) (string, error) {
	b, err := json.Marshal(struct {
		ExternalID   string `json:"external_id"`
		Type         string `json:"type"`
		FromWalletID string `json:"from_wallet_id"`
		ToWalletID   string `json:"to_wallet_id"`
		Amount       int64  `json:"amount"`
	}{
		ExternalID:   in.ExternalID,
		Type:         in.Type,
		FromWalletID: in.FromWalletID,
		ToWalletID:   in.ToWalletID,
		Amount:       in.Amount,
	})
	if err != nil {
		return "", err
	}
	// JSON payload is canonicalized enough for stable request matching in this service scope.
	return string(b), nil
}
