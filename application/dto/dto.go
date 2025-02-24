package dto

import "github.com/tomcyr/kryptonim-example/domain"

type RatesResponse struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

func NewRatesResponse(rate *domain.Rate) *RatesResponse {
	rateFloat, _ := rate.Rate.Round(rate.To.DecimalPoints).Float64()
	return &RatesResponse{
		From: rate.From.Symbol,
		To:   rate.To.Symbol,
		Rate: rateFloat,
	}
}

type ExchangeResponse struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func NewExchangeResponse(exchange *domain.Exchange) *ExchangeResponse {
	exchangeFloat, _ := exchange.Amount.Round(exchange.To.DecimalPoints).Float64()
	return &ExchangeResponse{
		From:   exchange.From.Symbol,
		To:     exchange.To.Symbol,
		Amount: exchangeFloat,
	}
}
