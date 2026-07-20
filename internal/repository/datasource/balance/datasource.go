package balance

import (
	"context"
	"fmt"

	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/service/balance"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type IDatasource interface {
	UpdateBalance(ctx context.Context, increment *balance.WalletIncrement) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (*balance.WalletIncrement, error)
}

type Datasource struct {
	pool *pgxpool.Pool
	log  *logrus.Logger
}

func NewDatasource(pool *pgxpool.Pool, log *logrus.Logger) IDatasource {
	return &Datasource{
		pool: pool,
		log:  log,
	}
}

func (d Datasource) UpdateBalance(ctx context.Context, increment *balance.WalletIncrement) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_transactions 
        (amount, wallet_id) 
   		VALUES ($1, $2)`, increment.Amount, increment.WalletID,
	)

	if err != nil {
		return fmt.Errorf("insert wallet transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (d Datasource) GetBalance(ctx context.Context, walletID uuid.UUID) (*balance.WalletIncrement, error) {
	// TODO implement me
	panic("implement me")
}
