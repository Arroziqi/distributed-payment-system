package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"transaction-service/internal/domain"
	"transaction-service/internal/repository"
)

type repoMock struct {
	processFn func(ctx context.Context, tx domain.Transaction, idemKey, requestHash string, expiresAt time.Time) (domain.Transaction, []byte, bool, error)
	listFn    func(ctx context.Context, limit, offset int) ([]domain.Transaction, error)
	getFn     func(ctx context.Context, id string) (domain.Transaction, error)
}

func (m repoMock) ProcessPayment(ctx context.Context, tx domain.Transaction, idemKey, requestHash string, expiresAt time.Time) (domain.Transaction, []byte, bool, error) {
	return m.processFn(ctx, tx, idemKey, requestHash, expiresAt)
}
func (m repoMock) List(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	return m.listFn(ctx, limit, offset)
}
func (m repoMock) GetByID(ctx context.Context, id string) (domain.Transaction, error) {
	return m.getFn(ctx, id)
}

type pubMock struct {
	published int
}

func (m *pubMock) Publish(_ context.Context, _ string, _ []byte) error {
	m.published++
	return nil
}

func TestProcessPaymentSuccess(t *testing.T) {
	pub := &pubMock{}
	repo := repoMock{
		processFn: func(_ context.Context, tx domain.Transaction, _ string, _ string, _ time.Time) (domain.Transaction, []byte, bool, error) {
			return tx, nil, false, nil
		},
		listFn: func(_ context.Context, _ int, _ int) ([]domain.Transaction, error) { return nil, nil },
		getFn:  func(_ context.Context, _ string) (domain.Transaction, error) { return domain.Transaction{}, nil },
	}
	uc := NewTransactionUsecase(repo, pub, 24*time.Hour)

	out, replayed, err := uc.ProcessPayment(context.Background(), ProcessPaymentInput{
		IdempotencyKey: "idem-1",
		Type:           "transfer",
		FromWalletID:   "w1",
		ToWalletID:     "w2",
		Amount:         100,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if replayed {
		t.Fatal("expected non-replayed response")
	}
	if out.ID == "" {
		t.Fatal("expected generated transaction id")
	}
	if pub.published != 1 {
		t.Fatalf("expected publish called once, got %d", pub.published)
	}
}

func TestProcessPaymentIdempotencyMismatch(t *testing.T) {
	pub := &pubMock{}
	repo := repoMock{
		processFn: func(_ context.Context, _ domain.Transaction, _ string, _ string, _ time.Time) (domain.Transaction, []byte, bool, error) {
			return domain.Transaction{}, nil, false, repository.ErrIdempotencyMismatch
		},
		listFn: func(_ context.Context, _ int, _ int) ([]domain.Transaction, error) { return nil, nil },
		getFn:  func(_ context.Context, _ string) (domain.Transaction, error) { return domain.Transaction{}, nil },
	}
	uc := NewTransactionUsecase(repo, pub, 24*time.Hour)

	_, _, err := uc.ProcessPayment(context.Background(), ProcessPaymentInput{
		IdempotencyKey: "idem-1",
		Type:           "transfer",
		FromWalletID:   "w1",
		ToWalletID:     "w2",
		Amount:         100,
	})
	if !errors.Is(err, ErrIdempotencyMismatch) {
		t.Fatalf("expected ErrIdempotencyMismatch, got %v", err)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	pub := &pubMock{}
	repo := repoMock{
		processFn: func(_ context.Context, tx domain.Transaction, _ string, _ string, _ time.Time) (domain.Transaction, []byte, bool, error) {
			return tx, nil, false, nil
		},
		listFn: func(_ context.Context, _ int, _ int) ([]domain.Transaction, error) { return nil, nil },
		getFn: func(_ context.Context, _ string) (domain.Transaction, error) {
			return domain.Transaction{}, repository.ErrNotFound
		},
	}
	uc := NewTransactionUsecase(repo, pub, 24*time.Hour)
	_, err := uc.GetByID(context.Background(), "missing")
	if !errors.Is(err, ErrTransactionNotFound) {
		t.Fatalf("expected ErrTransactionNotFound, got %v", err)
	}
}
