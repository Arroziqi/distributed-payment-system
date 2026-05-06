package postgres

import (
	"context"
	"errors"
	"time"

	"wallet-service/internal/domain"
	"wallet-service/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) WalletRepository {
	return WalletRepository{db: db}
}

func (r WalletRepository) Create(ctx context.Context, wallet domain.Wallet) error {
	const q = `
		INSERT INTO wallets (id, user_id, currency, available_balance, locked_balance, version, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, q, wallet.ID, wallet.UserID, wallet.Currency, wallet.AvailableBalance, wallet.LockedBalance, wallet.Version, wallet.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repository.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r WalletRepository) GetByUserID(ctx context.Context, userID string) (domain.Wallet, error) {
	const q = `
		SELECT id, user_id, currency, available_balance, locked_balance, version, updated_at
		FROM wallets WHERE user_id = $1
	`
	var w domain.Wallet
	err := r.db.QueryRow(ctx, q, userID).Scan(&w.ID, &w.UserID, &w.Currency, &w.AvailableBalance, &w.LockedBalance, &w.Version, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Wallet{}, repository.ErrNotFound
		}
		return domain.Wallet{}, err
	}
	return w, nil
}

func (r WalletRepository) Topup(ctx context.Context, userID string, amount int64, txnID string) (domain.Wallet, error) {
	return r.applySingleWalletMutation(ctx, userID, txnID, amount, "credit")
}

func (r WalletRepository) Withdraw(ctx context.Context, userID string, amount int64, txnID string) (domain.Wallet, error) {
	return r.applySingleWalletMutation(ctx, userID, txnID, -amount, "debit")
}

func (r WalletRepository) Transfer(ctx context.Context, fromUserID string, toUserID string, amount int64, txnID string) (domain.Wallet, domain.Wallet, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Wallet{}, domain.Wallet{}, err
	}
	defer tx.Rollback(ctx)

	from, err := r.getWalletByUserIDTx(ctx, tx, fromUserID)
	if err != nil {
		return domain.Wallet{}, domain.Wallet{}, err
	}
	to, err := r.getWalletByUserIDTx(ctx, tx, toUserID)
	if err != nil {
		return domain.Wallet{}, domain.Wallet{}, err
	}
	if from.AvailableBalance < amount {
		return domain.Wallet{}, domain.Wallet{}, repository.ErrInsufficientBalance
	}

	newFromBal := from.AvailableBalance - amount
	newToBal := to.AvailableBalance + amount
	now := time.Now().UTC()

	if err := r.updateWalletBalanceOptimisticTx(ctx, tx, from.ID, from.Version, newFromBal, now); err != nil {
		return domain.Wallet{}, domain.Wallet{}, err
	}
	if err := r.updateWalletBalanceOptimisticTx(ctx, tx, to.ID, to.Version, newToBal, now); err != nil {
		return domain.Wallet{}, domain.Wallet{}, err
	}

	if err := r.insertLedgerTx(ctx, tx, from.ID, txnID, "debit", amount, newFromBal, now); err != nil {
		return domain.Wallet{}, domain.Wallet{}, err
	}
	if err := r.insertLedgerTx(ctx, tx, to.ID, txnID, "credit", amount, newToBal, now); err != nil {
		return domain.Wallet{}, domain.Wallet{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Wallet{}, domain.Wallet{}, err
	}

	from.AvailableBalance = newFromBal
	from.Version++
	from.UpdatedAt = now
	to.AvailableBalance = newToBal
	to.Version++
	to.UpdatedAt = now
	return from, to, nil
}

func (r WalletRepository) applySingleWalletMutation(ctx context.Context, userID string, txnID string, delta int64, direction string) (domain.Wallet, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Wallet{}, err
	}
	defer tx.Rollback(ctx)

	w, err := r.getWalletByUserIDTx(ctx, tx, userID)
	if err != nil {
		return domain.Wallet{}, err
	}

	newBal := w.AvailableBalance + delta
	if newBal < 0 {
		return domain.Wallet{}, repository.ErrInsufficientBalance
	}
	now := time.Now().UTC()
	if err := r.updateWalletBalanceOptimisticTx(ctx, tx, w.ID, w.Version, newBal, now); err != nil {
		return domain.Wallet{}, err
	}
	amt := delta
	if amt < 0 {
		amt = -amt
	}
	if err := r.insertLedgerTx(ctx, tx, w.ID, txnID, direction, amt, newBal, now); err != nil {
		return domain.Wallet{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Wallet{}, err
	}

	w.AvailableBalance = newBal
	w.Version++
	w.UpdatedAt = now
	return w, nil
}

func (r WalletRepository) getWalletByUserIDTx(ctx context.Context, tx pgx.Tx, userID string) (domain.Wallet, error) {
	const q = `
		SELECT id, user_id, currency, available_balance, locked_balance, version, updated_at
		FROM wallets
		WHERE user_id = $1
	`
	var w domain.Wallet
	err := tx.QueryRow(ctx, q, userID).Scan(&w.ID, &w.UserID, &w.Currency, &w.AvailableBalance, &w.LockedBalance, &w.Version, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Wallet{}, repository.ErrNotFound
		}
		return domain.Wallet{}, err
	}
	return w, nil
}

func (r WalletRepository) updateWalletBalanceOptimisticTx(ctx context.Context, tx pgx.Tx, walletID string, expectedVersion int64, newBalance int64, updatedAt time.Time) error {
	const q = `
		UPDATE wallets
		SET available_balance = $1, version = version + 1, updated_at = $2
		WHERE id = $3 AND version = $4
	`
	cmd, err := tx.Exec(ctx, q, newBalance, updatedAt, walletID, expectedVersion)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return repository.ErrConcurrentUpdate
	}
	return nil
}

func (r WalletRepository) insertLedgerTx(ctx context.Context, tx pgx.Tx, walletID string, txnID string, direction string, amount int64, balanceAfter int64, createdAt time.Time) error {
	const q = `
		INSERT INTO ledger_entries (id, wallet_id, txn_id, direction, amount, balance_after, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(ctx, q, uuid.NewString(), walletID, txnID, direction, amount, balanceAfter, createdAt)
	return err
}
