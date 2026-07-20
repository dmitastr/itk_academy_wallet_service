package balance

import (
	"context"

	datasrouce "github.com/dmitastr/itk_academy_wallet_service/internal/repository/datasource/balance"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IWalletService interface {
	UpdateBalance(ctx context.Context, walletIncrement *WalletIncrement) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (*WalletBalance, error)
}

type WalletService struct {
	ds  datasrouce.IDatasource
	log *logrus.Logger
}

func NewWalletService(ds datasrouce.IDatasource, log *logrus.Logger) IWalletService {
	return &WalletService{ds: ds, log: log}
}

func (w WalletService) UpdateBalance(ctx context.Context, walletIncrement *WalletIncrement) error {
	if err := walletIncrement.CorrectedAmount(); err != nil {
		w.log.WithError(err).Errorln("Correction amount failed")
		return err
	}

	if err := w.ds.UpdateBalance(ctx, walletIncrement); err != nil {
		w.log.WithError(err).Errorln("UpdateBalance failed")
		return err
	}

	return nil
}

func (w WalletService) GetBalance(ctx context.Context, walletID uuid.UUID) (*WalletBalance, error) {
	// TODO implement me
	panic("implement me")
}
