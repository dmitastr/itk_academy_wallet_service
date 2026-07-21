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
	AddDeposit(ctx context.Context, increment *models.WalletIncrement) error
	AddWithdrawal(ctx context.Context, increment *models.WalletIncrement) error
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

// AddWithdrawal checks if there is sufficient amount and adds new withdrawal
func (d Datasource) AddWithdrawal(ctx context.Context, increment *models.WalletIncrement) error {
	d.log.WithField("walletID", increment.WalletID).WithField("amount", increment.Amount).Info("adding withdrawal to db")
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, increment.WalletID.String())
	if err != nil {
		return fmt.Errorf("acquire advisory lock: %w", err)
	}

	var currentBalance int64
	err = tx.QueryRow(ctx, `SELECT COALESCE(SUM(amount), 0) FROM wallet_transactions WHERE wallet_id = $1`,
		increment.WalletID).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("could not get current balance: %w", err)
	}

	if currentBalance < increment.GetAbsAmount() {
		return core.ErrInsufficientFunds
	}
	query := `INSERT INTO wallet_transactions (wallet_id, amount) VALUES ($1, $2)`
	if _, err = tx.Exec(ctx, query, increment.WalletID, increment.Amount); err != nil {
		return fmt.Errorf("failed to add transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// AddDeposit adds new deposit
func (d Datasource) AddDeposit(ctx context.Context, increment *models.WalletIncrement) error {
	d.log.WithField("walletID", increment.WalletID).WithField("amount", increment.Amount).Info("adding deposit to db")

	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO wallet_transactions (wallet_id, amount) VALUES ($1, $2)`
	if _, err = tx.Exec(ctx, query, increment.WalletID, increment.Amount); err != nil {
		return fmt.Errorf("failed to add transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// GetBalance get current balance as sum of all transactions
func (d Datasource) GetBalance(ctx context.Context, walletID uuid.UUID) (*models.WalletBalance, error) {
	d.log.WithField("walletID", walletID).Info("getting balance for wallet from db")

	rows, err := d.pool.Query(ctx, `
		SELECT wallet_id, COALESCE(SUM(amount), 0) AS balance, MAX(created_at) AS updated_at
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
