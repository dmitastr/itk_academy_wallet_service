package migrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmitastr/itk_academy_wallet_service/migrations"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
)

// Run applies all pending migrations inside a single transaction per step.
func Run(ctx context.Context, connString string, log *logrus.Logger) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(
		ctx,
		dbConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Info("Database connection pool established")

	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to create database migrations fs: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, connString)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	log.Info("Database migration succeeded")
	return pool, nil
}
