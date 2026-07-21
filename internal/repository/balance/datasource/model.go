package datasource

import (
	"time"

	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/models"
	"github.com/google/uuid"
)

type transactionRow struct {
	WalletID      uuid.UUID `db:"wallet_id"`
	Balance       int64     `db:"balance"`
	UpdatedAtLast time.Time `db:"updated_at"`
}

func (r transactionRow) toDomain() *models.WalletBalance {
	return &models.WalletBalance{
		WalletID:      r.WalletID,
		Balance:       r.Balance,
		UpdatedAtLast: r.UpdatedAtLast,
	}
}
