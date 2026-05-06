package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"auth-service/internal/domain"
	"auth-service/internal/repository"
)

type userRepoMock struct {
	usersByEmail map[string]domain.User
	usersByID    map[string]domain.User
	createErr    error
}

func (m *userRepoMock) Create(_ context.Context, user domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	if _, exists := m.usersByEmail[user.Email]; exists {
		return repository.ErrAlreadyExists
	}
	m.usersByEmail[user.Email] = user
	m.usersByID[user.ID] = user
	return nil
}

func (m *userRepoMock) GetByEmail(_ context.Context, email string) (domain.User, error) {
	if u, ok := m.usersByEmail[email]; ok {
		return u, nil
	}
	return domain.User{}, repository.ErrNotFound
}

func (m *userRepoMock) GetByID(_ context.Context, id string) (domain.User, error) {
	if u, ok := m.usersByID[id]; ok {
		return u, nil
	}
	return domain.User{}, repository.ErrNotFound
}

type refreshRepoMock struct {
	items map[string]string
}

func (m *refreshRepoMock) Store(_ context.Context, token string, userID string, _ time.Duration) error {
	m.items[token] = userID
	return nil
}

func (m *refreshRepoMock) Consume(_ context.Context, token string) (string, error) {
	v, ok := m.items[token]
	if !ok {
		return "", repository.ErrNotFound
	}
	delete(m.items, token)
	return v, nil
}

func (m *refreshRepoMock) Delete(_ context.Context, token string) error {
	delete(m.items, token)
	return nil
}

type hasherMock struct{}

func (hasherMock) Hash(plain string) (string, error) {
	return "hash:" + plain, nil
}

func (hasherMock) Compare(hashed string, plain string) error {
	if hashed != "hash:"+plain {
		return errors.New("mismatch")
	}
	return nil
}

type tokenManagerMock struct {
	refreshCounter int
}

func (m *tokenManagerMock) CreateAccessToken(userID string, _ string) (string, error) {
	return "access-" + userID, nil
}

func (m *tokenManagerMock) GenerateRefreshToken() (string, error) {
	m.refreshCounter++
	return fmt.Sprintf("refresh-token-%d", m.refreshCounter), nil
}

func TestRegisterSuccess(t *testing.T) {
	users := &userRepoMock{usersByEmail: map[string]domain.User{}, usersByID: map[string]domain.User{}}
	refresh := &refreshRepoMock{items: map[string]string{}}
	tokens := &tokenManagerMock{}
	uc := NewAuthUsecase(users, refresh, hasherMock{}, tokens, 15*time.Minute, 24*time.Hour)

	out, err := uc.Register(context.Background(), RegisterInput{
		Email:    "user@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Email != "user@example.com" {
		t.Fatalf("unexpected email: %s", out.Email)
	}
	if out.UserID == "" {
		t.Fatal("expected generated user id")
	}
}

func TestLoginSuccess(t *testing.T) {
	user := domain.User{
		ID:           "u-1",
		Email:        "user@example.com",
		PasswordHash: "hash:password123",
		Status:       "active",
	}
	users := &userRepoMock{
		usersByEmail: map[string]domain.User{"user@example.com": user},
		usersByID:    map[string]domain.User{"u-1": user},
	}
	refresh := &refreshRepoMock{items: map[string]string{}}
	tokens := &tokenManagerMock{}
	uc := NewAuthUsecase(users, refresh, hasherMock{}, tokens, 15*time.Minute, 24*time.Hour)

	out, err := uc.Login(context.Background(), LoginInput{
		Email:    "user@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatal("expected both access and refresh token")
	}
}

func TestRefreshRotatesToken(t *testing.T) {
	user := domain.User{
		ID:           "u-1",
		Email:        "user@example.com",
		PasswordHash: "hash:password123",
		Status:       "active",
	}
	users := &userRepoMock{
		usersByEmail: map[string]domain.User{"user@example.com": user},
		usersByID:    map[string]domain.User{"u-1": user},
	}
	refresh := &refreshRepoMock{items: map[string]string{"old-token": "u-1"}}
	tokens := &tokenManagerMock{}
	uc := NewAuthUsecase(users, refresh, hasherMock{}, tokens, 15*time.Minute, 24*time.Hour)

	out, err := uc.Refresh(context.Background(), RefreshInput{RefreshToken: "old-token"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.RefreshToken == "old-token" {
		t.Fatal("expected refresh token rotation")
	}
	if _, exists := refresh.items["old-token"]; exists {
		t.Fatal("old token should be consumed")
	}
	if _, exists := refresh.items[out.RefreshToken]; !exists {
		t.Fatal("new token should be stored")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	user := domain.User{
		ID:           "u-1",
		Email:        "user@example.com",
		PasswordHash: "hash:password123",
		Status:       "active",
	}
	users := &userRepoMock{
		usersByEmail: map[string]domain.User{"user@example.com": user},
		usersByID:    map[string]domain.User{"u-1": user},
	}
	refresh := &refreshRepoMock{items: map[string]string{}}
	tokens := &tokenManagerMock{}
	uc := NewAuthUsecase(users, refresh, hasherMock{}, tokens, 15*time.Minute, 24*time.Hour)

	_, err := uc.Login(context.Background(), LoginInput{
		Email:    "user@example.com",
		Password: "wrong",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
