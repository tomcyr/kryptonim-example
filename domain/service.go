package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/govalues/decimal"
	"strings"
)

type RatesService interface {
	GetRates(ctx context.Context, currencies []*Currency) ([]*Rate, error)
}

type ratesService struct {
	ratesRepository RatesRepository
}

func NewRatesService(ratesRepository RatesRepository) *ratesService {
	return &ratesService{
		ratesRepository: ratesRepository,
	}
}

func (s *ratesService) GetRates(ctx context.Context, currencies []*Currency) ([]*Rate, error) {
	baseRates, err := s.ratesRepository.GetRates(ctx, baseCurrency, currencies)
	if err != nil {
		return nil, fmt.Errorf("failed to get base rates: %w", err)
	}

	rates, err := s.generateCurrencyPairsWithRates(currencies, baseRates, baseCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to generate currency pairs with rates: %w", err)
	}

	return rates, nil
}

func (s *ratesService) generateCurrencyPairsWithRates(currencies []*Currency, baseRates map[Currency]float64, base *Currency) ([]*Rate, error) {
	pairs := make(map[string]decimal.Decimal)
	n := len(currencies)

	if n < 2 {
		return nil, errors.New("invalid number of currencies")
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				cur1, cur2 := currencies[i], currencies[j]
				_, exists := pairs[cur1.Symbol+"-"+cur2.Symbol]
				if exists {
					continue
				}

				rate1, exists1 := baseRates[*cur1]
				rate2, exists2 := baseRates[*cur2]

				if cur1.Symbol == base.Symbol {
					if exists2 {
						rate, _ := decimal.NewFromFloat64(rate2)
						pairs[cur1.Symbol+"-"+cur2.Symbol] = rate
						pairs[cur2.Symbol+"-"+cur1.Symbol], _ = rate.Inv()
					}
				} else if cur2.Symbol == base.Symbol {
					if exists1 {
						rate, _ := decimal.NewFromFloat64(rate1)
						pairs[cur1.Symbol+"-"+cur2.Symbol], _ = rate.Inv()
						pairs[cur2.Symbol+"-"+cur1.Symbol] = rate
					}
				} else if exists1 && exists2 {
					rateDec1, _ := decimal.NewFromFloat64(rate1)
					rateDec2, _ := decimal.NewFromFloat64(rate2)
					pairs[cur1.Symbol+"-"+cur2.Symbol], _ = rateDec2.Quo(rateDec1)
					pairs[cur2.Symbol+"-"+cur1.Symbol], _ = rateDec1.Quo(rateDec2)
				}
			}
		}
	}

	var rates []*Rate
	for pair, rate := range pairs {
		hyphenIndex := strings.Index(pair, "-")
		newRates, err := NewRate(pair[:hyphenIndex], pair[hyphenIndex+1:], rate)
		if err != nil {
			return nil, err
		}
		rates = append(rates, newRates)
	}

	return rates, nil
}

type ExchangeService interface {
	Exchange(ctx context.Context, exchange *Exchange) (*Exchange, error)
}

type exchangeService struct {
	ratesRepository RatesRepository
}

func NewExchangeService(ratesRepository RatesRepository) *exchangeService {
	return &exchangeService{
		ratesRepository: ratesRepository,
	}
}

func (s *exchangeService) Exchange(ctx context.Context, exchange *Exchange) (*Exchange, error) {
	currencies := []*Currency{exchange.From, exchange.To}
	baseRates, err := s.ratesRepository.GetRates(ctx, baseCurrency, currencies)
	if err != nil {
		return nil, fmt.Errorf("failed to get base rates: %w", err)
	}

	res, err := s.calculate(exchange, baseRates)
	if err != nil {
		return nil, fmt.Errorf("failed exchange amount: %w", err)
	}

	return res, nil
}

func (s *exchangeService) calculate(exchange *Exchange, baseRates map[Currency]float64) (*Exchange, error) {
	fromRate, existsFrom := baseRates[*exchange.From]
	if !existsFrom {
		return nil, errors.New("from base rate not exists")
	}
	toRate, existsTo := baseRates[*exchange.To]
	if !existsTo {
		return nil, errors.New("to base rate not exists")
	}

	fromRateDec, _ := decimal.NewFromFloat64(fromRate)
	toRateDec, _ := decimal.NewFromFloat64(toRate)

	baseAmount, err := exchange.Amount.Mul(fromRateDec)
	if err != nil {
		return nil, err
	}
	finalAmount, err := baseAmount.Quo(toRateDec)
	if err != nil {
		return nil, err
	}

	return NewExchangeFromCurrencies(exchange.From, exchange.To, finalAmount), nil
}
