package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionCompletedEvent struct {
	TransactionID string    `json:"transaction_id"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	Amount        int64     `json:"amount"`
	FromWalletID  string    `json:"from_wallet_id"`
	ToWalletID    string    `json:"to_wallet_id"`
	OccurredAt    time.Time `json:"occurred_at"`
}

type NotificationService struct {
	db *pgxpool.Pool
}

func NewNotificationService(db *pgxpool.Pool) NotificationService {
	return NotificationService{db: db}
}

func (s NotificationService) HandleTransactionCompleted(ctx context.Context, body []byte) error {
	var evt TransactionCompletedEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		return fmt.Errorf("decode event: %w", err)
	}

	subject := fmt.Sprintf("Transaction %s successful", evt.TransactionID)
	message := fmt.Sprintf(
		"Your transaction is completed.\nType: %s\nAmount: %d\nFrom Wallet: %s\nTo Wallet: %s\n",
		evt.Type, evt.Amount, evt.FromWalletID, evt.ToWalletID,
	)

	// Email simulation: in production replace with SMTP/provider integration.
	log.Printf("[EMAIL_SIMULATION] subject=%q body=%q", subject, message)

	payload, _ := json.Marshal(evt)
	_, err := s.db.Exec(ctx, `
		INSERT INTO notifications (id, event_id, channel, recipient, payload, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 'sent', NOW(), NOW())
	`, uuid.NewString(), uuid.NewString(), "email", evt.ToWalletID, payload)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}
	return nil
}
