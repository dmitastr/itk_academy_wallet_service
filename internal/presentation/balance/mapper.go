package balance

import (
	"time"

	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/models"
)

// ToModel converts from handler layer [IncrementRequest] to domain layer [models.WalletIncrement]
func (r IncrementRequest) ToModel() (models.WalletIncrement, error) {
	return models.WalletIncrement{
		WalletID:      r.WalletID,
		Amount:        r.Amount,
		OperationType: r.OperationType,
	}, nil
}

// ToResponse converts from domain layer [models.WalletBalance] to handler layer [BalanceResponse]
func ToResponse(b *models.WalletBalance) BalanceResponse {
	return BalanceResponse{
		WalletID:  b.WalletID.String(),
		Balance:   b.Balance,
		UpdatedAt: b.UpdatedAtLast.Format(time.RFC3339),
	}
}
