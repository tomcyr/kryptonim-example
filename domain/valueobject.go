package domain

import (
	"fmt"
	"github.com/govalues/decimal"
	"strings"
)

type Currency string

func NewCurrency(currency string) (*Currency, error) {
	if !isValidCurrency(currency) {
		return nil, fmt.Errorf("invalid argument for currency: %s", currency)
	}
	cur := Currency(strings.ToUpper(currency))

	return &cur, nil
}

func isValidCurrency(currency string) bool {
	if len(currency) < 3 || len(currency) > 5 {
		return false
	}

	return true
}

func (c *Currency) String() string {
	return string(*c)
}

type Rates struct {
	From *Currency       `json:"from"`
	To   *Currency       `json:"to"`
	Rate decimal.Decimal `json:"rate"`
}

func NewRates(from, to string, rate decimal.Decimal, scale int) (*Rates, error) {
	f, err := NewCurrency(from)
	if err != nil {
		return nil, err
	}
	t, err := NewCurrency(to)
	if err != nil {
		return nil, err
	}

	return &Rates{
		From: f,
		To:   t,
		Rate: rate.Round(scale),
	}, nil
}
