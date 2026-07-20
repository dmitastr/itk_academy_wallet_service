package balance

import (
	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/models"
)

func (r IncrementRequest) ToModel() (models.WalletIncrement, error) {
	return models.WalletIncrement{
		WalletID:      r.WalletID,
		Amount:        r.Amount,
		OperationType: r.OperationType,
	}, nil
}
