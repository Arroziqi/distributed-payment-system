package repository

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrConcurrentUpdate    = errors.New("concurrent update")
)
