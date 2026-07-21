package service

import (
	"errors"
	"testing"

	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/models"
	mocklogger "github.com/dmitastr/itk_academy_wallet_service/internal/mocks/mock-logger"
	"github.com/dmitastr/itk_academy_wallet_service/internal/repository/datasource/balance/mocks"
	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestWalletService_GetBalance(t *testing.T) {
	walletID := uuid.New()

	tests := []struct {
		name            string
		walletID        uuid.UUID
		mockSetup       func(m *mocks.IDatasource)
		expectedBalance *models.WalletBalance
		expectedError   error
	}{
		{
			name:     "success",
			walletID: walletID,
			mockSetup: func(m *mocks.IDatasource) {
				m.EXPECT().
					GetBalance(mock.Anything, walletID).
					Return(&models.WalletBalance{WalletID: walletID, Amount: 1}, nil)
			},
			expectedBalance: &models.WalletBalance{WalletID: walletID, Amount: 1},
			expectedError:   nil,
		},
		{
			name:     "wallet not found",
			walletID: walletID,
			mockSetup: func(m *mocks.IDatasource) {
				m.EXPECT().
					GetBalance(mock.Anything, walletID).
					Return(nil, core.ErrWalletNotFound)
			},
			expectedBalance: nil,
			expectedError:   core.ErrWalletNotFound,
		},
		{
			name:     "db error",
			walletID: walletID,
			mockSetup: func(m *mocks.IDatasource) {
				m.EXPECT().
					GetBalance(mock.Anything, walletID).
					Return(nil, errors.New("db error"))
			},
			expectedBalance: nil,
			expectedError:   errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewIDatasource(t)
			tt.mockSetup(mockService)

			s := NewWalletService(mockService, mocklogger.NewTestLogger())
			balance, err := s.GetBalance(t.Context(), tt.walletID)

			assert.Equal(t, balance, tt.expectedBalance)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestWalletService_AddDeposit(t *testing.T) {
	walletID := uuid.New()

	tests := []struct {
		name          string
		increment     *models.WalletIncrement
		mockSetup     func(m *mocks.IDatasource)
		expectedError error
	}{
		{
			name:      "success deposit",
			increment: &models.WalletIncrement{Amount: 1, WalletID: walletID, OperationType: core.OperationTypeDeposit},
			mockSetup: func(m *mocks.IDatasource) {
				m.EXPECT().
					AddDeposit(mock.Anything, mock.MatchedBy(func(incr *models.WalletIncrement) bool {
						return incr.Amount == 1 && incr.WalletID == walletID && incr.OperationType == core.OperationTypeDeposit
					})).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "success withdraw",
			increment: &models.WalletIncrement{Amount: 1, WalletID: walletID, OperationType: core.OperationTypeWithdraw},
			mockSetup: func(m *mocks.IDatasource) {
				m.EXPECT().
					AddWithdrawal(mock.Anything, mock.MatchedBy(func(incr *models.WalletIncrement) bool {
						return incr.Amount == -1 && incr.WalletID == walletID && incr.OperationType == core.OperationTypeWithdraw
					})).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:      "deposit negative amount",
			increment: &models.WalletIncrement{Amount: -1, WalletID: walletID, OperationType: core.OperationTypeDeposit},
			mockSetup: func(m *mocks.IDatasource) {
				m.EXPECT().
					AddDeposit(mock.Anything, mock.MatchedBy(func(incr *models.WalletIncrement) bool {
						return incr.Amount == 1 && incr.WalletID == walletID && incr.OperationType == core.OperationTypeDeposit
					})).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "unknown operation",
			increment:     &models.WalletIncrement{Amount: 1, WalletID: walletID, OperationType: core.OperationType("invalid operation type")},
			mockSetup:     func(m *mocks.IDatasource) {},
			expectedError: core.ErrInvalidOperationType,
		},
		{
			name:      "db error",
			increment: &models.WalletIncrement{Amount: 1, WalletID: walletID, OperationType: core.OperationTypeDeposit},
			mockSetup: func(m *mocks.IDatasource) {
				m.EXPECT().
					AddDeposit(mock.Anything, mock.Anything).
					Return(errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewIDatasource(t)
			tt.mockSetup(mockService)

			s := NewWalletService(mockService, mocklogger.NewTestLogger())
			err := s.AddDeposit(t.Context(), tt.increment)

			assert.Equal(t, tt.expectedError, err)
		})
	}
}
