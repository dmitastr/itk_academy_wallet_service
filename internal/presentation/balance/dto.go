package balance

import (
	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/google/uuid"
)

type IncrementRequest struct {
	WalletID      uuid.UUID          `json:"valletId"`
	OperationType core.OperationType `json:"operationType"`
	Amount        int64              `json:"amount"`
}

type BalanceResponse struct {
	WalletID  string `json:"valletId"`
	Amount    int64  `json:"amount"`     // строка, чтобы не терять точность
	UpdatedAt string `json:"updated_at"` // RFC3339
}

type SuccessResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
