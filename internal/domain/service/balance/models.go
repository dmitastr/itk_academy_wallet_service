package balance

import (
	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/google/uuid"
)

type WalletIncrement struct {
	WalletID      uuid.UUID
	Amount        int64
	OperationType core.OperationType
}

func (i WalletIncrement) CorrectedAmount() error {
	var sign int64 = 1
	if i.Amount < 0 {
		i.Amount *= -1
	}

	switch i.OperationType {
	case core.OperationTypeWithdraw:
		sign = -1
	case core.OperationTypeDeposit:
		sign = 1
	default:
		return core.ErrInvalidOperationType
	}
	i.Amount *= sign
	return nil

}

type WalletBalance struct {
	WalletID uuid.UUID
	Amount   int64
}
