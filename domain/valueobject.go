package domain

import (
	"fmt"
	"github.com/govalues/decimal"
	"strings"
)

const (
	defaultFiatDecimalPoints int = 6
)

type Currency struct {
	Symbol        string
	DecimalPoints int
}

func NewCurrency(symbol string) (*Currency, error) {
	if !isValidSymbol(symbol) {
		return nil, fmt.Errorf("invalid argument for currency: %s", symbol)
	}
	cur := Currency{
		Symbol:        strings.ToUpper(symbol),
		DecimalPoints: defaultFiatDecimalPoints,
	}

	return &cur, nil
}

func isValidSymbol(symbol string) bool {
	if len(symbol) < 3 || len(symbol) > 5 {
		return false
	}

	return true
}

func (c *Currency) String() string {
	return c.Symbol
}

func (c *Currency) SetDecimalPoints(decimalPoints int) {
	c.DecimalPoints = decimalPoints
}

type Rate struct {
	From *Currency
	To   *Currency
	Rate decimal.Decimal
}

func NewRate(from, to string, rate decimal.Decimal) (*Rate, error) {
	f, err := NewCurrency(from)
	if err != nil {
		return nil, err
	}
	t, err := NewCurrency(to)
	if err != nil {
		return nil, err
	}

	return &Rate{
		From: f,
		To:   t,
		Rate: rate,
	}, nil
}

type Exchange struct {
	From   *Currency
	To     *Currency
	Amount decimal.Decimal
}

func NewExchange(from, to string, amount float64) (*Exchange, error) {
	fromCurr, err := NewCurrency(from)
	if err != nil {
		return nil, err
	}
	toCurr, err := NewCurrency(to)
	if err != nil {
		return nil, err
	}
	amountDec, err := decimal.NewFromFloat64(amount)
	if err != nil {
		return nil, err
	}

	return &Exchange{
		From:   fromCurr,
		To:     toCurr,
		Amount: amountDec,
	}, nil
}

func NewExchangeFromCurrencies(from, to *Currency, amount decimal.Decimal) *Exchange {
	return &Exchange{
		From:   from,
		To:     to,
		Amount: amount,
	}
}
