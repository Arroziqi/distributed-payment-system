package repository

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrIdempotencyMismatch = errors.New("idempotency request mismatch")
	ErrIdempotencyBusy     = errors.New("idempotency key is being processed")
	ErrConcurrentUpdate    = errors.New("concurrent update")
)
