package domain

import (
	"context"
)

var baseCurrency, _ = NewCurrency("USD")

type RatesRepository interface {
	GetRates(ctx context.Context, baseCurrency *Currency, currencies []*Currency) (map[Currency]float64, error)
}
