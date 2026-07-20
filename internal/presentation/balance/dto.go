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

type BalanceResoponse struct {
	WalletID uuid.UUID `json:"valletId"`
	Balance  int64     `json:"balance"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
