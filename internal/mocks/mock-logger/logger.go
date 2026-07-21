package mock_logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

func NewTestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(io.Discard) // логи никуда не выводятся, консоль теста чистая
	return logger
}
