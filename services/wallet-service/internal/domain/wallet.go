package domain

import "time"

type Wallet struct {
	ID               string
	UserID           string
	Currency         string
	AvailableBalance int64
	LockedBalance    int64
	Version          int64
	UpdatedAt        time.Time
}

type LedgerEntry struct {
	ID           string
	WalletID      string
	TxnID         string
	Direction     string
	Amount        int64
	BalanceAfter  int64
	CreatedAt     time.Time
}
