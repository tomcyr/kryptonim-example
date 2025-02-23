package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tomcyr/kryptonim-example/domain"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	latestEndpoint = "latest.json"
	apiTimeout     = 5 * time.Second
)

type openExchangeRatesRepository struct {
	appID   string
	baseURL string
}

func NewOpenExchangeRatesRepository(appID, baseURL string) *openExchangeRatesRepository {
	return &openExchangeRatesRepository{
		appID:   appID,
		baseURL: baseURL,
	}
}

type ratesResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

type apiError struct {
	IsError     bool   `json:"error"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (e apiError) Error() string {
	return fmt.Sprintf("%v: %v", e.Message, e.Description)
}

func (r *openExchangeRatesRepository) GetRates(ctx context.Context, baseCurrency domain.Currency, currencies []*domain.Currency) (map[string]float64, error) {
	ctx, cancel := context.WithTimeout(ctx, apiTimeout)
	defer cancel()

	args := make(map[string]string)
	args["base"] = baseCurrency.String()

	if currencies != nil {
		symbols := make([]string, len(currencies))
		for i, currency := range currencies {
			symbols[i] = currency.String()
		}
		args["symbols"] = strings.Join(symbols, ",")
	}

	data, err := r.apiCall(latestEndpoint, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get rates: %w", err)
	}

	var ratesRes ratesResponse
	err = json.Unmarshal(data, &ratesRes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal rates response: %w", err)
	}

	return ratesRes.Rates, nil
}

func (r *openExchangeRatesRepository) apiCall(endpoint string, args map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%v/%v?app_id=%v", r.baseURL, endpoint, r.appID)
	for k := range args {
		url = fmt.Sprintf("%v&%v=%v", url, k, args[k])
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var e apiError
		err = json.Unmarshal(body, &e)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal error during api call: %w", err)
		}
		return nil, fmt.Errorf("response failure during api call: %w", e)
	}

	return body, nil
}
