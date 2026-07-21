package testhelpers

import (
	"context"
	"time"

	"github.com/dmitastr/itk_academy_wallet_service/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	*config.DBConfig
	ConnStr string
}

// CreatePostgresContainer creates container for running tests on database
func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	user := "postgres"
	password := "postgres"
	dbName := "wallet-db"
	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")
	connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnStr:           connStr,
		DBConfig: &config.DBConfig{
			Host: host,
			Port: port.Port(),
			User: user,
			Pass: password,
			Name: dbName,
		},
	}, nil
}
