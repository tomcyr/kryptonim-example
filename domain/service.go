package domain

import (
	"context"
	"fmt"
	"github.com/govalues/decimal"
	"strings"
)

type RatesService interface {
	GetRates(ctx context.Context, currencies []*Currency) ([]*Rates, error)
}

type ratesService struct {
	ratesRepository RatesRepository
}

func NewRatesService(ratesRepository RatesRepository) *ratesService {
	return &ratesService{
		ratesRepository: ratesRepository,
	}
}

func (s *ratesService) GetRates(ctx context.Context, currencies []*Currency) ([]*Rates, error) {
	baseRates, err := s.ratesRepository.GetRates(ctx, baseCurrency, currencies)
	if err != nil {
		return nil, err
	}

	rates, err := s.generateCurrencyPairsWithRates(currencies, baseRates, baseCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to generate currency pairs with rates: %w", err)
	}

	return rates, nil
}

func (s *ratesService) generateCurrencyPairsWithRates(currencies []*Currency, baseRates map[string]float64, base Currency) ([]*Rates, error) {
	pairs := make(map[string]decimal.Decimal)
	var rates []*Rates
	n := len(currencies)

	if n < 2 {
		return rates, nil
	}

	for i := 0; i < len(currencies); i++ {
		for j := 0; j < len(currencies); j++ {
			if i != j {
				cur1, cur2 := currencies[i].String(), currencies[j].String()
				_, exists := pairs[cur1+"-"+cur2]
				if exists {
					continue
				}

				rate1, exists1 := baseRates[cur1]
				rate2, exists2 := baseRates[cur2]

				if cur1 == base.String() {
					if exists2 {
						rate, _ := decimal.NewFromFloat64(rate2)
						pairs[cur1+"-"+cur2] = rate
						pairs[cur2+"-"+cur1], _ = rate.Inv()
					}
				} else if cur2 == base.String() {
					if exists1 {
						rate, _ := decimal.NewFromFloat64(rate1)
						pairs[cur1+"-"+cur2], _ = rate.Inv()
						pairs[cur2+"-"+cur1] = rate
					}
				} else if exists1 && exists2 {
					rateDec1, _ := decimal.NewFromFloat64(rate1)
					rateDec2, _ := decimal.NewFromFloat64(rate2)
					pairs[cur1+"-"+cur2], _ = rateDec2.Quo(rateDec1)
					pairs[cur2+"-"+cur1], _ = rateDec1.Quo(rateDec2)
				}
			}
		}
	}

	for pair, rate := range pairs {
		hyphenIndex := strings.Index(pair, "-")
		newRates, err := NewRates(pair[:hyphenIndex], pair[hyphenIndex+1:], rate, defaultFiatDecimalPoints)
		if err != nil {
			return nil, err
		}
		rates = append(rates, newRates)
	}

	return rates, nil
}
