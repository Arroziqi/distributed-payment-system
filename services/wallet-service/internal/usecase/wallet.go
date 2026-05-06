package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"wallet-service/internal/domain"
	"wallet-service/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput        = errors.New("invalid input")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrConcurrentUpdate    = errors.New("concurrent update")
	ErrWalletExists        = errors.New("wallet already exists")
)

type WalletUsecase struct {
	walletRepo repository.WalletRepository
	cacheRepo  repository.BalanceCacheRepository
	cacheTTL   time.Duration
}

func NewWalletUsecase(walletRepo repository.WalletRepository, cacheRepo repository.BalanceCacheRepository, cacheTTL time.Duration) *WalletUsecase {
	return &WalletUsecase{walletRepo: walletRepo, cacheRepo: cacheRepo, cacheTTL: cacheTTL}
}

func (u *WalletUsecase) CreateWallet(ctx context.Context, userID string, currency string) (domain.Wallet, error) {
	if userID == "" {
		return domain.Wallet{}, ErrInvalidInput
	}
	if currency == "" {
		currency = "USD"
	}
	wallet := domain.Wallet{
		ID:               uuid.NewString(),
		UserID:           userID,
		Currency:         currency,
		AvailableBalance: 0,
		LockedBalance:    0,
		Version:          0,
		UpdatedAt:        time.Now().UTC(),
	}
	if err := u.walletRepo.Create(ctx, wallet); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return domain.Wallet{}, ErrWalletExists
		}
		return domain.Wallet{}, fmt.Errorf("create wallet: %w", err)
	}
	_ = u.cacheRepo.Set(ctx, userID, 0, u.cacheTTL)
	return wallet, nil
}

func (u *WalletUsecase) Topup(ctx context.Context, userID string, amount int64) (domain.Wallet, error) {
	if userID == "" || amount <= 0 {
		return domain.Wallet{}, ErrInvalidInput
	}
	w, err := u.walletRepo.Topup(ctx, userID, amount, uuid.NewString())
	if err != nil {
		return domain.Wallet{}, mapRepoErr(err)
	}
	_ = u.cacheRepo.Set(ctx, userID, w.AvailableBalance, u.cacheTTL)
	return w, nil
}

func (u *WalletUsecase) Withdraw(ctx context.Context, userID string, amount int64) (domain.Wallet, error) {
	if userID == "" || amount <= 0 {
		return domain.Wallet{}, ErrInvalidInput
	}
	w, err := u.walletRepo.Withdraw(ctx, userID, amount, uuid.NewString())
	if err != nil {
		return domain.Wallet{}, mapRepoErr(err)
	}
	_ = u.cacheRepo.Set(ctx, userID, w.AvailableBalance, u.cacheTTL)
	return w, nil
}

func (u *WalletUsecase) Transfer(ctx context.Context, fromUserID string, toUserID string, amount int64) (domain.Wallet, domain.Wallet, error) {
	if fromUserID == "" || toUserID == "" || fromUserID == toUserID || amount <= 0 {
		return domain.Wallet{}, domain.Wallet{}, ErrInvalidInput
	}
	fromWallet, toWallet, err := u.walletRepo.Transfer(ctx, fromUserID, toUserID, amount, uuid.NewString())
	if err != nil {
		return domain.Wallet{}, domain.Wallet{}, mapRepoErr(err)
	}
	_ = u.cacheRepo.Set(ctx, fromUserID, fromWallet.AvailableBalance, u.cacheTTL)
	_ = u.cacheRepo.Set(ctx, toUserID, toWallet.AvailableBalance, u.cacheTTL)
	return fromWallet, toWallet, nil
}

func (u *WalletUsecase) BalanceInquiry(ctx context.Context, userID string) (domain.Wallet, error) {
	if userID == "" {
		return domain.Wallet{}, ErrInvalidInput
	}

	w, err := u.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return domain.Wallet{}, mapRepoErr(err)
	}

	_ = u.cacheRepo.Set(ctx, userID, w.AvailableBalance, u.cacheTTL)

	return w, nil
}

func mapRepoErr(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return ErrWalletNotFound
	case errors.Is(err, repository.ErrInsufficientBalance):
		return ErrInsufficientBalance
	case errors.Is(err, repository.ErrConcurrentUpdate):
		return ErrConcurrentUpdate
	case errors.Is(err, repository.ErrAlreadyExists):
		return ErrWalletExists
	default:
		return err
	}
}
