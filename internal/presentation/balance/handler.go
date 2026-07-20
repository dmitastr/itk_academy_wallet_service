package balance

import (
	"net/http"

	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	. "github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IWalletHandlers interface {
	UpdateBalance(ctx *gin.Context)
	GetBalance(ctx *gin.Context)
}

type WalletHandlers struct {
	service IWalletService
	log     *logrus.Logger
}

func NewWalletHandlers(service IWalletService, logger *logrus.Logger) IWalletHandlers {
	return &WalletHandlers{service: service, log: logger}
}

func (w WalletHandlers) UpdateBalance(ctx *gin.Context) {
	var request IncrementRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		w.log.WithError(err).Error("failed to bind request body")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	model, err := request.ToModel()
	if err != nil {
		w.log.WithError(err).Error("failed to convert request to model")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
	if err := w.service.UpdateBalance(ctx, &model); err != nil {
		w.HandleServiceError(ctx, err)
		w.log.WithError(err).Error("failed to update balance")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	ctx.JSON(http.StatusOK, SuccessResponse{Data: request})
}

func (w WalletHandlers) GetBalance(ctx *gin.Context) {
	walletID := ctx.Param("wallet_id")
	if walletID == "" {
		w.log.Error(core.ErrMissingWalletID.Error())
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: core.ErrMissingWalletID.Error()})
		return
	}

	walletUUID, err := uuid.FromBytes([]byte(walletID))
	if err != nil {
		w.log.WithError(err).Error("error converting wallet uuid")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: core.ErrInvalidWalletID.Error()})
		return
	}
	balance, err := w.service.GetBalance(ctx, walletUUID)
	if err != nil {
		w.HandleServiceError(ctx, err)
		w.log.WithError(err).Error("failed to get balance")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	ctx.JSON(http.StatusOK, SuccessResponse{Data: ToResponse(balance)})
}
