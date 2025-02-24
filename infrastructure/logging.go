package infrastructure

import (
	"context"
	"github.com/tomcyr/kryptonim-example/domain"
	"go.uber.org/zap"
)

type loggingRatesRepository struct {
	wrapped domain.RatesRepository
	logger  *zap.Logger
}

func NewLoggingRatesRepository(wrapped domain.RatesRepository, logger *zap.Logger) *loggingRatesRepository {
	return &loggingRatesRepository{
		wrapped: wrapped,
		logger:  logger,
	}
}

func (l *loggingRatesRepository) GetRates(ctx context.Context, baseCurrency *domain.Currency, currencies []*domain.Currency) (map[domain.Currency]float64, error) {
	l.logger.Debug("get rates from repository", zap.String("baseCurrency", baseCurrency.String()), zap.Reflect("currencies", currencies))
	res, err := l.wrapped.GetRates(ctx, baseCurrency, currencies)
	if err != nil {
		l.logger.Error("failed to get rates from repository", zap.Error(err))

		return res, err
	}

	resLog := make(map[string]float64, len(res))
	for currency, f := range res {
		resLog[currency.String()] = f
	}
	l.logger.Debug("get rates from repository successfully", zap.String("baseCurrency", baseCurrency.String()), zap.Reflect("currencies", currencies), zap.Any("result", resLog))

	return res, err
}
