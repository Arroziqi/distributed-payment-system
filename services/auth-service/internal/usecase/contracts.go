package usecase

import "context"

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hashed string, plain string) error
}

type TokenManager interface {
	CreateAccessToken(userID string, email string) (string, error)
	GenerateRefreshToken() (string, error)
}

type UserFinder interface {
	GetUserByID(ctx context.Context, userID string) (UserView, error)
}

type UserView struct {
	ID    string `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
