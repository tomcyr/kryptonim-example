package domain

import (
	"context"
)

const (
	baseCurrency             Currency = "USD"
	defaultFiatDecimalPoints int      = 6
)

type RatesRepository interface {
	GetRates(ctx context.Context, baseCurrency Currency, currencies []*Currency) (map[string]float64, error)
}
