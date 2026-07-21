package balance

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/models"
	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/service/mocks"
	mocklogger "github.com/dmitastr/itk_academy_wallet_service/internal/mocks/mock-logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWalletHandlers_UpdateBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	walletID := uuid.New()

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(m *mocks.IWalletService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "success",
			requestBody: fmt.Sprintf(`{"amount": 100, "operationType": "DEPOSIT", "valletId": "%s"}`, walletID),
			mockSetup: func(m *mocks.IWalletService) {
				m.EXPECT().
					AddDeposit(mock.Anything, mock.MatchedBy(func(model *models.WalletIncrement) bool {
						return model.WalletID == walletID && model.Amount == 100
					})).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "insufficient funds",
			requestBody: fmt.Sprintf(`{"amount": 100, "operationType": "DEPOSIT", "valletId": "%s"}`, walletID),
			mockSetup: func(m *mocks.IWalletService) {
				m.EXPECT().
					AddDeposit(mock.Anything, mock.Anything).
					Return(core.ErrInsufficientFunds)
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "invalid request body",
			requestBody:    `{"amount": ""}`,
			mockSetup:      func(m *mocks.IWalletService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid wallet id",
			requestBody:    fmt.Sprintf(`{"amount": 100, "operationType": "DEPOSIT", "valletId": "%s"}`, "walletID"),
			mockSetup:      func(m *mocks.IWalletService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid operation type",
			requestBody: fmt.Sprintf(`{"amount": 100, "operationType": "BALANCE", "valletId": "%s"}`, walletID),
			mockSetup: func(m *mocks.IWalletService) {
				m.EXPECT().
					AddDeposit(mock.Anything, mock.Anything).
					Return(core.ErrInvalidOperationType)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewIWalletService(t)
			tt.mockSetup(mockService)

			h := NewWalletHandlers(mockService, mocklogger.NewTestLogger())
			router := gin.New()
			router.POST("/wallet", h.UpdateBalance)

			req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestWalletHandlers_GetBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	walletID := uuid.New()

	tests := []struct {
		name           string
		walletID       string
		mockSetup      func(m *mocks.IWalletService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "success",
			walletID: walletID.String(),
			mockSetup: func(m *mocks.IWalletService) {
				m.EXPECT().
					GetBalance(mock.Anything, walletID).
					Return(&models.WalletBalance{WalletID: walletID, Amount: 1}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf(`{"data": {"balance": 1, "valletId": "%s"}}`, walletID),
		},
		{
			name:     "wallet not found",
			walletID: walletID.String(),
			mockSetup: func(m *mocks.IWalletService) {
				m.EXPECT().
					GetBalance(mock.Anything, mock.Anything).
					Return(nil, core.ErrWalletNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid wallet id",
			walletID:       "invalid-uuid",
			mockSetup:      func(m *mocks.IWalletService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := mocks.NewIWalletService(t)
			tt.mockSetup(mockService)

			h := NewWalletHandlers(mockService, mocklogger.NewTestLogger())
			router := gin.New()
			router.GET("/wallets/:walletID", h.GetBalance)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/wallets/%s", tt.walletID), nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
