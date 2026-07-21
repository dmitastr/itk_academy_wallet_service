// repository/txretry/retry.go
package txretry

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	MaxRetries int
	BaseDelay  time.Duration
	MaxDelay   time.Duration
}

func DefaultConfig() Config {
	return Config{
		MaxRetries: 3,
		BaseDelay:  50 * time.Millisecond,
		MaxDelay:   500 * time.Millisecond,
	}
}

// WithRetry выполняет fn в транзакции, повторяя при транзиентных ошибках БД
func WithRetry(ctx context.Context, pool *pgxpool.Pool, cfg Config, isoLevel pgx.TxIsoLevel, fn func(tx pgx.Tx) error) error {
	var lastErr error

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := backoffDelay(cfg, attempt)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		err := runOnce(ctx, pool, isoLevel, fn)
		if err == nil {
			return nil
		}

		if !isRetryable(err) {
			return err // не транзиентная ошибка — сразу возвращаем, без повторов
		}

		lastErr = err
	}

	return fmt.Errorf("transaction failed after %d retries: %w", cfg.MaxRetries, lastErr)
}

func runOnce(ctx context.Context, pool *pgxpool.Pool, isoLevel pgx.TxIsoLevel, fn func(tx pgx.Tx) error) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func isRetryable(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "40001": // serialization_failure
			return true
		case "40P01": // deadlock_detected
			return true
		}
	}
	return false
}

func backoffDelay(cfg Config, attempt int) time.Duration {
	delay := cfg.BaseDelay * time.Duration(1<<uint(attempt-1)) // экспоненциальный рост
	if delay > cfg.MaxDelay {
		delay = cfg.MaxDelay
	}
	// jitter, чтобы конкурентные ретраи не столкнулись снова одновременно
	jitter := time.Duration(rand.Int63n(int64(delay) / 2))
	return delay + jitter
}
