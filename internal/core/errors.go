package core

import "errors"

var (
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrInvalidOperationType = errors.New("invalid operation type")
	ErrMissingWalletID      = errors.New("wallet_id is required")
	ErrInvalidWalletID      = errors.New("invalid wallet_id")
)
