package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dmitastr/itk_academy_wallet_service/internal/config"
	balance2 "github.com/dmitastr/itk_academy_wallet_service/internal/domain/service/balance"
	"github.com/dmitastr/itk_academy_wallet_service/internal/presentation/balance"
	balance3 "github.com/dmitastr/itk_academy_wallet_service/internal/repository/datasource/balance"
	"github.com/dmitastr/itk_academy_wallet_service/internal/repository/migrations"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WalletServiceApp struct {
	server *http.Server
	cfg    config.ConfigProvider
	log    *logrus.Logger
}

func NewWalletServiceApp(ctx context.Context, cfg config.ConfigProvider, log *logrus.Logger) (*WalletServiceApp, error) {
	app := &WalletServiceApp{cfg: cfg, log: log}
	router := gin.Default()

	pool, err := migrations.Run(ctx, cfg.GetDBConfig(), log)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	ds := balance3.NewDatasource(pool, log)
	walletService := balance2.NewWalletService(ds, log)
	walletHandlers := balance.NewWalletHandlers(walletService, log)

	app.RegisterRoutes(router, walletHandlers)

	return app, nil
}

func (app *WalletServiceApp) Start() error {
	if err := app.server.ListenAndServe(); err != nil {
		return fmt.Errorf("start server failed: %s", err)
	}
	return nil
}

func (app *WalletServiceApp) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := app.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server failed: %s", err)
	}
	return nil
}

func (app *WalletServiceApp) RegisterRoutes(router *gin.Engine, walletHandlers balance.IWalletHandlers) {
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	apiPath := router.Group("/api/v1/wallet")
	apiPath.GET(`/wallets/:walletID`, walletHandlers.GetBalance)
	apiPath.POST(`/wallet`, walletHandlers.UpdateBalance)
}
