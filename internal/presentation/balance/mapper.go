package balance

import (
	"time"

	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/models"
)

func (r IncrementRequest) ToModel() (models.WalletIncrement, error) {
	return models.WalletIncrement{
		WalletID:      r.WalletID,
		Amount:        r.Amount,
		OperationType: r.OperationType,
	}, nil
}

func ToResponse(b *models.WalletBalance) BalanceResponse {
	return BalanceResponse{
		WalletID:  b.WalletID.String(),
		Balance:   b.Amount,
		UpdatedAt: b.UpdatedAtLast.Format(time.RFC3339),
	}
}
