package balance

import (
	"errors"
	"net/http"

	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/gin-gonic/gin"
)

var errStatusMap = map[error]int{
	core.ErrInsufficientFunds:    http.StatusUnprocessableEntity,
	core.ErrWalletNotFound:       http.StatusNotFound,
	core.ErrInvalidOperationType: http.StatusBadRequest,
}

func (w WalletHandlers) HandleServiceError(c *gin.Context, err error) {
	for sentinel, status := range errStatusMap {
		if errors.Is(err, sentinel) {
			w.log.WithError(err).Error("internal server error")
			c.JSON(status, ErrorResponse{Error: sentinel.Error()})
			return
		}
	}

	w.log.WithError(err).Error("unhandled service error")
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "internal server error",
	})
}
