package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"transaction-service/internal/domain"
	"transaction-service/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return TransactionRepository{db: db}
}

func (r TransactionRepository) ProcessPayment(ctx context.Context, txData domain.Transaction, idemKey, requestHash string, expiresAt time.Time) (domain.Transaction, []byte, bool, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return domain.Transaction{}, nil, false, err
	}
	defer tx.Rollback(ctx)

	record, inserted, err := r.upsertAndLockIdempotency(ctx, tx, idemKey, requestHash, expiresAt)
	if err != nil {
		return domain.Transaction{}, nil, false, err
	}
	if !inserted {
		if record.RequestHash != requestHash {
			return domain.Transaction{}, nil, false, repository.ErrIdempotencyMismatch
		}
		if record.Status == "completed" && len(record.Response) > 0 {
			var prev domain.Transaction
			if err := json.Unmarshal(record.Response, &prev); err != nil {
				return domain.Transaction{}, nil, false, err
			}
			if err := tx.Commit(ctx); err != nil {
				return domain.Transaction{}, nil, false, err
			}
			return prev, record.Response, true, nil
		}
		if record.Status == "processing" {
			return domain.Transaction{}, nil, false, repository.ErrIdempotencyBusy
		}
	}

	const insertTxn = `
		INSERT INTO transactions (id, external_id, type, from_wallet_id, to_wallet_id, amount, status, version, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`
	_, err = tx.Exec(
		ctx,
		insertTxn,
		txData.ID,
		nullIfEmpty(txData.ExternalID),
		txData.Type,
		txData.FromWalletID,
		txData.ToWalletID,
		txData.Amount,
		txData.Status,
		txData.Version,
		txData.CreatedAt,
		txData.UpdatedAt,
	)
	if err != nil {
		return domain.Transaction{}, nil, false, err
	}

	responseSnapshot, err := json.Marshal(txData)
	if err != nil {
		return domain.Transaction{}, nil, false, err
	}

	const updateIdem = `
		UPDATE idempotency_keys
		SET status='completed', response_snapshot=$1
		WHERE key=$2
	`
	if _, err := tx.Exec(ctx, updateIdem, responseSnapshot, idemKey); err != nil {
		return domain.Transaction{}, nil, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Transaction{}, nil, false, err
	}
	return txData, responseSnapshot, false, nil
}

func (r TransactionRepository) List(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	const q = `
		SELECT id, COALESCE(external_id, ''), type, from_wallet_id, to_wallet_id, amount, status, version, created_at, updated_at
		FROM transactions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.ExternalID, &t.Type, &t.FromWalletID, &t.ToWalletID, &t.Amount, &t.Status, &t.Version, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, rows.Err()
}

func (r TransactionRepository) GetByID(ctx context.Context, id string) (domain.Transaction, error) {
	const q = `
		SELECT id, COALESCE(external_id, ''), type, from_wallet_id, to_wallet_id, amount, status, version, created_at, updated_at
		FROM transactions WHERE id=$1
	`
	var t domain.Transaction
	err := r.db.QueryRow(ctx, q, id).Scan(&t.ID, &t.ExternalID, &t.Type, &t.FromWalletID, &t.ToWalletID, &t.Amount, &t.Status, &t.Version, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Transaction{}, repository.ErrNotFound
		}
		return domain.Transaction{}, err
	}
	return t, nil
}

func (r TransactionRepository) upsertAndLockIdempotency(ctx context.Context, tx pgx.Tx, key, reqHash string, expiresAt time.Time) (domain.IdempotencyRecord, bool, error) {
	const insert = `
		INSERT INTO idempotency_keys (id, key, request_hash, status, expires_at, created_at)
		VALUES ($1, $2, $3, 'processing', $4, NOW())
		ON CONFLICT (key) DO NOTHING
	`
	tag, err := tx.Exec(ctx, insert, uuid.NewString(), key, reqHash, expiresAt)
	if err != nil {
		return domain.IdempotencyRecord{}, false, err
	}
	inserted := tag.RowsAffected() == 1

	const selectForUpdate = `
		SELECT key, request_hash, status, COALESCE(response_snapshot, '{}'::jsonb), expires_at
		FROM idempotency_keys
		WHERE key = $1
		FOR UPDATE
	`
	var rec domain.IdempotencyRecord
	if err := tx.QueryRow(ctx, selectForUpdate, key).Scan(&rec.Key, &rec.RequestHash, &rec.Status, &rec.Response, &rec.ExpiresAt); err != nil {
		return domain.IdempotencyRecord{}, false, err
	}
	return rec, inserted, nil
}

func nullIfEmpty(v string) any {
	if v == "" {
		return nil
	}
	return v
}
