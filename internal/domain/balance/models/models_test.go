package models

import (
	"testing"

	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestWalletIncrement_CorrectedAmount(t *testing.T) {
	tests := []struct {
		name          string
		amount        int64
		operationType core.OperationType
		expected      int64
	}{
		{
			name:          "positive deposit",
			amount:        10,
			operationType: core.OperationTypeDeposit,
			expected:      10,
		},
		{
			name:          "negative deposit",
			amount:        -10,
			operationType: core.OperationTypeDeposit,
			expected:      10,
		},
		{
			name:          "positive withdrawal",
			amount:        10,
			operationType: core.OperationTypeWithdraw,
			expected:      -10,
		},
		{
			name:          "negative withdrawal",
			amount:        -10,
			operationType: core.OperationTypeWithdraw,
			expected:      -10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := WalletIncrement{Amount: tt.amount, OperationType: tt.operationType}
			err := w.CorrectedAmount()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, w.Amount)
		})
	}
}
