package balance

import (
	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/google/uuid"
)

type IncrementRequest struct {
	WalletID      uuid.UUID          `json:"valletId" binding:"required,uuid4"`
	OperationType core.OperationType `json:"operationType" binding:"required"`
	Amount        int64              `json:"amount" binding:"required"`
}

type BalanceResponse struct {
	WalletID  string `json:"valletId"`
	Balance   int64  `json:"balance"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
