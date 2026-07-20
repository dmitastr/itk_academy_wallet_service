package balance

import (
	"net/http"

	. "github.com/dmitastr/itk_academy_wallet_service/internal/domain/service/balance"
	"github.com/gin-gonic/gin"
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

	ctx.JSON(http.StatusNoContent, request)

}

func (w WalletHandlers) GetBalance(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}
