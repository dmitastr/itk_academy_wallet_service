package core

import "errors"

var (
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrInvalidOperationType = errors.New("invalid operation type")
)
