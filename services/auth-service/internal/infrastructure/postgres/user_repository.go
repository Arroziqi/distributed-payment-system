package postgres

import (
	"context"
	"errors"
	"fmt"

	"auth-service/internal/domain"
	"auth-service/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Create(ctx context.Context, user domain.User) error {
	const q = `
		INSERT INTO users (id, name, email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, q, user.ID, user.Name, user.Email, user.PasswordHash, user.Status, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: %s", repository.ErrAlreadyExists, "duplicate email")
		}
		return err
	}
	return nil
}

func (r UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	const q = `
		SELECT id, name, email, password_hash, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var u domain.User
	err := r.db.QueryRow(ctx, q, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, repository.ErrNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	const q = `
		SELECT id, name, email, password_hash, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u domain.User
	err := r.db.QueryRow(ctx, q, id).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, repository.ErrNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}
func (r UserRepository) Update(ctx context.Context, user domain.User) error {
	const q = `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, status = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(ctx, q, user.Name, user.Email, user.PasswordHash, user.Status, user.UpdatedAt, user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: %s", repository.ErrAlreadyExists, "duplicate email")
		}
		return err
	}
	return nil
}
