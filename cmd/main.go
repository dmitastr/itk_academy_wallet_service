package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/dmitastr/itk_academy_wallet_service/internal/app"
	"github.com/dmitastr/itk_academy_wallet_service/internal/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	log := logrus.New()
	logrus.SetLevel(logrus.DebugLevel)
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer func() {
		stop()
	}()

	walletApp, err := app.NewWalletServiceApp(ctx, cfg, log)
	if err != nil {
		panic(err)
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := walletApp.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("app run failed: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		return walletApp.Stop(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
