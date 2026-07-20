package balance

import (
	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/service/balance"
)

func (r IncrementRequest) ToModel() (balance.WalletIncrement, error) {
	return balance.WalletIncrement{
		WalletID:      r.WalletID,
		Amount:        r.Amount,
		OperationType: r.OperationType,
	}, nil
}
