package domain

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u User) IsActive() bool {
	return u.Status == "active"
}
