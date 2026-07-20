package balance

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type IDatasource interface {
	UpdateBalance(ctx context.Context, increment *models.WalletIncrement) error
	GetBalance(ctx context.Context, walletID uuid.UUID) (*models.WalletBalance, error)
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

func (d Datasource) UpdateBalance(ctx context.Context, increment *models.WalletIncrement) error {
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

func (d Datasource) GetBalance(ctx context.Context, walletID uuid.UUID) (*models.WalletBalance, error) {
	rows, err := d.pool.Query(ctx, `
		SELECT wallet_id, COALESCE(SUM(amount), 0) AS balance, MAX(updated_at) AS updated_at
		FROM wallet_transactions
		WHERE wallet_id = $1
		GROUP BY wallet_id`, walletID)

	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}

	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[transactionRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.ErrWalletNotFound
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	return row.toDomain(), nil
}
