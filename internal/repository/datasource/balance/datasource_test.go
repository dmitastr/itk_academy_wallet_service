package balance

import (
	"context"
	"log"
	"testing"

	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	mocklogger "github.com/dmitastr/itk_academy_wallet_service/internal/mocks/mock-logger"
	"github.com/dmitastr/itk_academy_wallet_service/internal/mocks/testhelpers"
	"github.com/dmitastr/itk_academy_wallet_service/internal/repository/migrations"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestNewDatasource(t *testing.T) {
	suite.Run(t, new(DatasourceTestSuite))
}

type DatasourceTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  IDatasource
	ctx         context.Context
	pool        *pgxpool.Pool
}

func (suite *DatasourceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	logger := mocklogger.NewTestLogger()

	pool, err := migrations.Run(suite.ctx, pgContainer.ConnStr, logger)
	suite.Require().NoError(err)
	repository := NewDatasource(pool, logger)
	if err != nil {
		log.Fatal(err)
	}
	suite.repository = repository
	suite.pool = pool
}

func (suite *DatasourceTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.T().Log("Error while terminating PG container", err)
	}
}

func (suite *DatasourceTestSuite) TestDatasource_GetBalance() {
	walletID := uuid.New()
	walletIDSum := uuid.New()

	_, err := suite.pool.Exec(suite.ctx, `
        INSERT INTO wallet_transactions (wallet_id, amount)
        VALUES ($1, $2), ($3, $4), ($5, $6)
    `, walletID, 100, walletIDSum, 10, walletIDSum, 10)
	require.NoError(suite.T(), err)

	tests := []struct {
		name           string
		walletID       uuid.UUID
		expectedAmount int64
		expectedError  error
	}{
		{
			name:           "success single entry",
			walletID:       walletID,
			expectedAmount: 100,
			expectedError:  nil,
		},
		{
			name:           "success multiple entry",
			walletID:       walletIDSum,
			expectedAmount: 20,
			expectedError:  nil,
		},
		{
			name:          "no transactions",
			walletID:      uuid.New(),
			expectedError: core.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			balance, err := suite.repository.GetBalance(suite.ctx, tt.walletID)
			if tt.expectedError != nil {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), tt.expectedError.Error(), err.Error())
				return
			}
			assert.Equal(suite.T(), tt.expectedAmount, balance.Amount)
		})
	}

}
