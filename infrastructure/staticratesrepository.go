package infrastructure

import (
	"context"
	"errors"
	"github.com/tomcyr/kryptonim-example/domain"
)

var rates = []rate{
	{CryptoCurrency: "BEER", DecimalPlaces: 18, Rate: 0.00002461, BaseCurrency: "USD"},
	{CryptoCurrency: "FLOKI", DecimalPlaces: 18, Rate: 0.0001428, BaseCurrency: "USD"},
	{CryptoCurrency: "GATE", DecimalPlaces: 18, Rate: 6.87, BaseCurrency: "USD"},
	{CryptoCurrency: "USDT", DecimalPlaces: 6, Rate: 0.999, BaseCurrency: "USD"},
	{CryptoCurrency: "WBTC", DecimalPlaces: 8, Rate: 57037.22, BaseCurrency: "USD"},
}

type rate struct {
	CryptoCurrency string
	DecimalPlaces  int
	Rate           float64
	BaseCurrency   string
}

type staticRatesRepository struct{}

func NewStaticRatesRepository() *staticRatesRepository {
	return &staticRatesRepository{}
}

func (r *staticRatesRepository) GetRates(_ context.Context, baseCurrency *domain.Currency, currencies []*domain.Currency) (map[domain.Currency]float64, error) {
	result := make(map[domain.Currency]float64)
	for _, currency := range currencies {
		for _, rat := range rates {
			if rat.CryptoCurrency == currency.Symbol && rat.BaseCurrency == baseCurrency.Symbol {
				currency.SetDecimalPoints(rat.DecimalPlaces)
				result[*currency] = rat.Rate
				break
			}
		}
	}

	if len(result) < 2 {
		return nil, errors.New("at least one currency not exist in rate repository")
	}

	return result, nil
}
