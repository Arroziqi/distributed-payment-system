package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"auth-service/internal/domain"
	"auth-service/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInactiveUser        = errors.New("inactive user")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidInput        = errors.New("invalid input")
)

type AuthUsecase struct {
	users      repository.UserRepository
	refresh    repository.RefreshTokenRepository
	hasher     PasswordHasher
	token      TokenManager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthUsecase(
	users repository.UserRepository,
	refresh repository.RefreshTokenRepository,
	hasher PasswordHasher,
	token TokenManager,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *AuthUsecase {
	return &AuthUsecase{
		users:      users,
		refresh:    refresh,
		hasher:     hasher,
		token:      token,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *AuthUsecase) Register(ctx context.Context, in RegisterInput) (RegisterOutput, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || len(in.Password) < 8 {
		return RegisterOutput{}, ErrInvalidInput
	}

	_, err := u.users.GetByEmail(ctx, email)
	if err == nil {
		return RegisterOutput{}, ErrEmailAlreadyExists
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return RegisterOutput{}, fmt.Errorf("check existing user: %w", err)
	}

	passwordHash, err := u.hasher.Hash(in.Password)
	if err != nil {
		return RegisterOutput{}, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	user := domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: passwordHash,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return RegisterOutput{}, ErrEmailAlreadyExists
		}
		return RegisterOutput{}, fmt.Errorf("create user: %w", err)
	}

	return RegisterOutput{
		UserID:    user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	UserId       string `json:"user_id"`
}

func (u *AuthUsecase) Login(ctx context.Context, in LoginInput) (AuthTokens, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || in.Password == "" {
		return AuthTokens{}, ErrInvalidInput
	}

	user, err := u.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return AuthTokens{}, ErrInvalidCredentials
		}
		return AuthTokens{}, fmt.Errorf("get user by email: %w", err)
	}
	if !user.IsActive() {
		return AuthTokens{}, ErrInactiveUser
	}

	if err := u.hasher.Compare(user.PasswordHash, in.Password); err != nil {
		return AuthTokens{}, ErrInvalidCredentials
	}

	return u.issueTokens(ctx, user.ID, user.Email)
}

type RefreshInput struct {
	RefreshToken string
}

func (u *AuthUsecase) Refresh(ctx context.Context, in RefreshInput) (AuthTokens, error) {
	if strings.TrimSpace(in.RefreshToken) == "" {
		return AuthTokens{}, ErrInvalidInput
	}

	userID, err := u.refresh.Consume(ctx, in.RefreshToken)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return AuthTokens{}, ErrInvalidRefreshToken
		}
		return AuthTokens{}, fmt.Errorf("consume refresh token: %w", err)
	}

	user, err := u.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return AuthTokens{}, ErrInvalidRefreshToken
		}
		return AuthTokens{}, fmt.Errorf("get user by id: %w", err)
	}
	if !user.IsActive() {
		return AuthTokens{}, ErrInactiveUser
	}

	return u.issueTokens(ctx, user.ID, user.Email)
}

type LogoutInput struct {
	RefreshToken string
}

func (u *AuthUsecase) Logout(ctx context.Context, in LogoutInput) error {
	if strings.TrimSpace(in.RefreshToken) == "" {
		return ErrInvalidInput
	}
	if err := u.refresh.Delete(ctx, in.RefreshToken); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

func (u *AuthUsecase) issueTokens(ctx context.Context, userID string, email string) (AuthTokens, error) {
	accessToken, err := u.token.CreateAccessToken(userID, email)
	if err != nil {
		return AuthTokens{}, fmt.Errorf("create access token: %w", err)
	}

	refreshToken, err := u.token.GenerateRefreshToken()
	if err != nil {
		return AuthTokens{}, fmt.Errorf("create refresh token: %w", err)
	}

	if err := u.refresh.Store(ctx, refreshToken, userID, u.refreshTTL); err != nil {
		return AuthTokens{}, fmt.Errorf("store refresh token: %w", err)
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(u.accessTTL.Seconds()),
		UserId:       userID,
	}, nil
}
