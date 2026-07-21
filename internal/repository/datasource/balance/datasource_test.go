package balance

import (
	"context"
	"log"
	"sync"
	"testing"

	"github.com/dmitastr/itk_academy_wallet_service/internal/core"
	"github.com/dmitastr/itk_academy_wallet_service/internal/domain/balance/models"
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

func (s *DatasourceTestSuite) SetupSuite() {
	s.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(s.ctx)
	if err != nil {
		log.Fatal(err)
	}
	s.pgContainer = pgContainer
	logger := mocklogger.NewTestLogger()

	pool, err := migrations.Run(s.ctx, pgContainer.ConnStr, logger)
	s.Require().NoError(err)
	repository := NewDatasource(pool, logger)
	if err != nil {
		log.Fatal(err)
	}
	s.repository = repository
	s.pool = pool
}

func (s *DatasourceTestSuite) TearDownSuite() {
	if err := s.pgContainer.Terminate(s.ctx); err != nil {
		s.T().Log("Error while terminating PG container", err)
	}
}

func (s *DatasourceTestSuite) SetupTest() {
	ctx := context.Background()
	_, err := s.pool.Exec(ctx, `TRUNCATE wallet_transactions`)
	s.Require().NoError(err)
}

func (s *DatasourceTestSuite) countTransactions(walletID uuid.UUID) (int, int64) {
	s.T().Helper()
	var count int
	var sum int64
	err := s.pool.QueryRow(context.Background(), `
        SELECT COUNT(*), COALESCE(SUM(amount), 0)
        FROM wallet_transactions
        WHERE wallet_id = $1
    `, walletID).Scan(&count, &sum)
	s.Require().NoError(err)
	return count, sum
}

func (s *DatasourceTestSuite) TestDatasource_GetBalance() {
	walletID := uuid.New()
	walletIDSum := uuid.New()

	_, err := s.pool.Exec(s.ctx, `
        INSERT INTO wallet_transactions (wallet_id, amount)
        VALUES ($1, $2), ($3, $4), ($5, $6)
    `, walletID, 100, walletIDSum, 10, walletIDSum, 10)
	require.NoError(s.T(), err)

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
		s.Run(tt.name, func() {
			balance, err := s.repository.GetBalance(s.ctx, tt.walletID)
			if tt.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tt.expectedError.Error(), err.Error())
				return
			}
			assert.Equal(s.T(), tt.expectedAmount, balance.Amount)
		})
	}

}

func (s *DatasourceTestSuite) TestDatasource_AddDeposit_Success() {
	walletID := uuid.New()
	increment := &models.WalletIncrement{
		WalletID: walletID,
		Amount:   100,
	}

	err := s.repository.AddTransaction(context.Background(), increment)
	s.Require().NoError(err)

	count, sum := s.countTransactions(walletID)
	s.Equal(1, count)
	s.True(sum == increment.Amount)

}

func (s *DatasourceTestSuite) TestAddDeposit_ConcurrentInserts_NoneLost() {
	walletID := uuid.New()
	const concurrency = 50
	var wg sync.WaitGroup
	errCh := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.repository.AddTransaction(context.Background(), &models.WalletIncrement{
				WalletID: walletID,
				Amount:   1,
			})
			errCh <- err
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		s.NoError(err)
	}

	count, _ := s.countTransactions(walletID)
	s.Equal(concurrency, count, "all concurrent deposits must be persisted, none lost")
}

func (s *DatasourceTestSuite) TestDatasource_AddWithdrawal_Success() {
	var amount int64 = 100
	walletID := uuid.New()
	increment := &models.WalletIncrement{
		WalletID: walletID,
		Amount:   -1 * amount,
	}

	_, err := s.pool.Exec(context.Background(), `INSERT INTO wallet_transactions (wallet_id, amount) VALUES ($1, $2)`, walletID, amount)
	s.Require().NoError(err)

	err = s.repository.AddWithdrawal(context.Background(), increment)
	s.Require().NoError(err)

	count, sum := s.countTransactions(walletID)
	s.Equal(2, count)
	s.True(sum == 0)
}

func (s *DatasourceTestSuite) TestDatasource_AddWithdrawal_InsufficientAmount() {
	var amountDeposit int64 = 50
	var amountWithdraw int64 = 100
	walletID := uuid.New()
	increment := &models.WalletIncrement{
		WalletID: walletID,
		Amount:   -1 * amountWithdraw,
	}

	_, err := s.pool.Exec(context.Background(), `INSERT INTO wallet_transactions (wallet_id, amount) VALUES ($1, $2)`, walletID, amountDeposit)
	s.Require().NoError(err)

	err = s.repository.AddWithdrawal(context.Background(), increment)
	s.Require().Error(err)
	s.Equal(err, core.ErrInsufficientFunds)
}

func (s *DatasourceTestSuite) TestDatasource_AddWithdrawal_walletIDNotFound() {
	var amountWithdraw int64 = 100
	walletID := uuid.New()
	increment := &models.WalletIncrement{
		WalletID: walletID,
		Amount:   -1 * amountWithdraw,
	}

	err := s.repository.AddWithdrawal(context.Background(), increment)
	s.Require().Error(err)
	s.Equal(core.ErrInsufficientFunds, err)
}

func (s *DatasourceTestSuite) TestAddWithdrawal_ConcurrentInserts_NoneLost() {
	walletID := uuid.New()

	var amount int64 = 100
	_, err := s.pool.Exec(context.Background(), `INSERT INTO wallet_transactions (wallet_id, amount) VALUES ($1, $2)`, walletID, amount)
	s.Require().NoError(err)
	const concurrency = 50
	var wg sync.WaitGroup
	errCh := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.repository.AddWithdrawal(context.Background(), &models.WalletIncrement{
				WalletID: walletID,
				Amount:   -1,
			})
			errCh <- err
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		s.NoError(err)
	}

	count, sum := s.countTransactions(walletID)
	s.Equal(concurrency+1, count, "all concurrent withdrawals must be persisted, none lost")
	s.True(sum == (amount - concurrency))
}
