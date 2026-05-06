package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"wallet-service/internal/domain"
	"wallet-service/internal/repository"
)

type walletRepoMock struct {
	byUser map[string]domain.Wallet
}

func (m *walletRepoMock) Create(_ context.Context, wallet domain.Wallet) error {
	if _, ok := m.byUser[wallet.UserID]; ok {
		return repository.ErrAlreadyExists
	}
	m.byUser[wallet.UserID] = wallet
	return nil
}

func (m *walletRepoMock) GetByUserID(_ context.Context, userID string) (domain.Wallet, error) {
	w, ok := m.byUser[userID]
	if !ok {
		return domain.Wallet{}, repository.ErrNotFound
	}
	return w, nil
}

func (m *walletRepoMock) Topup(_ context.Context, userID string, amount int64, _ string) (domain.Wallet, error) {
	w, ok := m.byUser[userID]
	if !ok {
		return domain.Wallet{}, repository.ErrNotFound
	}
	w.AvailableBalance += amount
	w.Version++
	m.byUser[userID] = w
	return w, nil
}

func (m *walletRepoMock) Withdraw(_ context.Context, userID string, amount int64, _ string) (domain.Wallet, error) {
	w, ok := m.byUser[userID]
	if !ok {
		return domain.Wallet{}, repository.ErrNotFound
	}
	if w.AvailableBalance < amount {
		return domain.Wallet{}, repository.ErrInsufficientBalance
	}
	w.AvailableBalance -= amount
	w.Version++
	m.byUser[userID] = w
	return w, nil
}

func (m *walletRepoMock) Transfer(_ context.Context, fromUserID string, toUserID string, amount int64, _ string) (domain.Wallet, domain.Wallet, error) {
	from, ok := m.byUser[fromUserID]
	if !ok {
		return domain.Wallet{}, domain.Wallet{}, repository.ErrNotFound
	}
	to, ok := m.byUser[toUserID]
	if !ok {
		return domain.Wallet{}, domain.Wallet{}, repository.ErrNotFound
	}
	if from.AvailableBalance < amount {
		return domain.Wallet{}, domain.Wallet{}, repository.ErrInsufficientBalance
	}
	from.AvailableBalance -= amount
	to.AvailableBalance += amount
	from.Version++
	to.Version++
	m.byUser[fromUserID] = from
	m.byUser[toUserID] = to
	return from, to, nil
}

type cacheRepoMock struct {
	cache map[string]int64
}

func (m *cacheRepoMock) Get(_ context.Context, userID string) (int64, bool, error) {
	v, ok := m.cache[userID]
	return v, ok, nil
}
func (m *cacheRepoMock) Set(_ context.Context, userID string, balance int64, _ time.Duration) error {
	m.cache[userID] = balance
	return nil
}
func (m *cacheRepoMock) Invalidate(_ context.Context, userID string) error {
	delete(m.cache, userID)
	return nil
}

func TestCreateWallet(t *testing.T) {
	repo := &walletRepoMock{byUser: map[string]domain.Wallet{}}
	cache := &cacheRepoMock{cache: map[string]int64{}}
	uc := NewWalletUsecase(repo, cache, 5*time.Minute)

	w, err := uc.CreateWallet(context.Background(), "user-1", "USD")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if w.UserID != "user-1" {
		t.Fatalf("unexpected user id: %s", w.UserID)
	}
}

func TestTopupAndBalance(t *testing.T) {
	repo := &walletRepoMock{
		byUser: map[string]domain.Wallet{
			"user-1": {UserID: "user-1", AvailableBalance: 100},
		},
	}
	cache := &cacheRepoMock{cache: map[string]int64{}}
	uc := NewWalletUsecase(repo, cache, 5*time.Minute)

	w, err := uc.Topup(context.Background(), "user-1", 50)
	if err != nil {
		t.Fatalf("topup failed: %v", err)
	}
	if w.AvailableBalance != 150 {
		t.Fatalf("expected 150, got %d", w.AvailableBalance)
	}
	balance, err := uc.BalanceInquiry(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("balance inquiry failed: %v", err)
	}
	if balance != 150 {
		t.Fatalf("expected 150, got %d", balance)
	}
}

func TestWithdrawInsufficient(t *testing.T) {
	repo := &walletRepoMock{
		byUser: map[string]domain.Wallet{
			"user-1": {UserID: "user-1", AvailableBalance: 20},
		},
	}
	cache := &cacheRepoMock{cache: map[string]int64{}}
	uc := NewWalletUsecase(repo, cache, 5*time.Minute)

	_, err := uc.Withdraw(context.Background(), "user-1", 100)
	if !errors.Is(err, ErrInsufficientBalance) {
		t.Fatalf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestTransfer(t *testing.T) {
	repo := &walletRepoMock{
		byUser: map[string]domain.Wallet{
			"sender":   {UserID: "sender", AvailableBalance: 200},
			"receiver": {UserID: "receiver", AvailableBalance: 20},
		},
	}
	cache := &cacheRepoMock{cache: map[string]int64{}}
	uc := NewWalletUsecase(repo, cache, 5*time.Minute)

	from, to, err := uc.Transfer(context.Background(), "sender", "receiver", 50)
	if err != nil {
		t.Fatalf("transfer failed: %v", err)
	}
	if from.AvailableBalance != 150 || to.AvailableBalance != 70 {
		t.Fatalf("unexpected balances from=%d to=%d", from.AvailableBalance, to.AvailableBalance)
	}
}
