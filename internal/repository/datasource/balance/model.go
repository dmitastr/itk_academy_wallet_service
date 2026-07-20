package balance

import (
	"time"

	"github.com/google/uuid"
)

type transactionRow struct {
	WalletID  uuid.UUID `db:"wallet_id"`
	Amount    int64     `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}
