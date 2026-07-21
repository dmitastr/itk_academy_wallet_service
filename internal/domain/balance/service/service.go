package service

import (
	"context"

	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/models"
	datasrouce "github.com/dmitastr/itk_academy_wallet_service/internal/repository/datasource/balance"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IWalletService interface {
	AddTransaction(ctx context.Context, walletIncrement *models.WalletIncrement) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (*models.WalletBalance, error)
}

type WalletService struct {
	ds  datasrouce.IDatasource
	log *logrus.Logger
}

func NewWalletService(ds datasrouce.IDatasource, log *logrus.Logger) IWalletService {
	return &WalletService{ds: ds, log: log}
}

// AddTransaction handles correcting transaction amount and selecting appropriate db method
func (w WalletService) AddTransaction(ctx context.Context, walletIncrement *models.WalletIncrement) error {
	if err := walletIncrement.CorrectedAmount(); err != nil {
		w.log.WithError(err).Errorln("Correction amount failed")
		return err
	}

	var err error
	switch walletIncrement.OperationType {
	case core.OperationTypeDeposit:
		err = w.ds.AddDeposit(ctx, walletIncrement)
	case core.OperationTypeWithdraw:
		err = w.ds.AddWithdrawal(ctx, walletIncrement)
	default:
		return core.ErrInvalidOperationType

	}
	return err
}

// GetBalance fetches balance of a wallet from db
func (w WalletService) GetBalance(ctx context.Context, walletID uuid.UUID) (*models.WalletBalance, error) {
	w.log.WithField("wallet_id", walletID).Debugln("GetBalance")
	balance, err := w.ds.GetBalance(ctx, walletID)
	if err != nil {
		w.log.WithError(err).Errorln("GetBalance failed")
		return nil, err
	}

	return balance, nil
}
