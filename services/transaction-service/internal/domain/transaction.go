package domain

import "time"

type Transaction struct {
	ID           string
	ExternalID   string
	Type         string
	FromWalletID string
	ToWalletID   string
	Amount       int64
	Status       string
	Version      int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type IdempotencyRecord struct {
	Key         string
	RequestHash string
	Status      string
	Response    []byte
	ExpiresAt   time.Time
}
