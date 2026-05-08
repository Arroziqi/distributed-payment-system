package domain

import "time"

type Wallet struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Currency         string    `json:"currency"`
	AvailableBalance int64     `json:"available_balance"`
	LockedBalance    int64     `json:"locked_balance"`
	Version          int64     `json:"version"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type LedgerEntry struct {
	ID           string    `json:"id"`
	WalletID     string    `json:"wallet_id"`
	TxnID        string    `json:"txn_id"`
	Direction    string    `json:"direction"`
	Amount       int64     `json:"amount"`
	BalanceAfter int64     `json:"balance_after"`
	CreatedAt    time.Time `json:"created_at"`
}

