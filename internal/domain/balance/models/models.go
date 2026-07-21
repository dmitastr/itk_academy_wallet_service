package models

import (
	"time"

	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/google/uuid"
)

type WalletIncrement struct {
	WalletID      uuid.UUID
	Amount        int64
	OperationType core.OperationType
}

// CorrectedAmount fix any amount anomalies and make it positive for deposit and negative for withdrawal
func (i *WalletIncrement) CorrectedAmount() error {
	var sign int64 = 1
	switch i.OperationType {
	case core.OperationTypeWithdraw:
		sign = -1
	case core.OperationTypeDeposit:
		sign = 1
	default:
		return core.ErrInvalidOperationType
	}

	i.Amount = sign * abs(i.Amount)
	return nil
}

func (i *WalletIncrement) GetAbsAmount() int64 {
	return abs(i.Amount)
}

type WalletBalance struct {
	WalletID      uuid.UUID
	Balance       int64
	UpdatedAtLast time.Time
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
