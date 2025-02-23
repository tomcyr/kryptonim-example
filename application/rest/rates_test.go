package rest

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tomcyr/kryptonim-example/domain"
	"github.com/tomcyr/kryptonim-example/mocks"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRates_GetRatesCurrencies_Success(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	ctx, r := gin.CreateTestContext(responseRecorder)
	repositoryMock := mocks.NewMockRatesRepository(t)
	service := domain.NewRatesService(repositoryMock)

	eur, _ := domain.NewCurrency("EUR")
	gbp, _ := domain.NewCurrency("GBP")
	currencies := []*domain.Currency{eur, gbp}

	baseRates := make(map[string]float64, 2)
	baseRates["EUR"] = 1
	baseRates["GBP"] = 2

	baseCurrency, _ := domain.NewCurrency("USD")

	repositoryMock.On("GetRates", mock.Anything, *baseCurrency, currencies).Return(
		baseRates,
		nil,
	)

	ratesHandler := NewRatesHandler(service, zap.NewNop())
	req, _ := http.NewRequest("GET", "/rates?currencies=EUR,GBP", nil)
	r.GET("/rates", ratesHandler.GetRatesCurrencies)
	ctx.Request = req

	r.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	repositoryMock.AssertExpectations(t)
	var response []domain.Rates
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}

func TestRates_GetRatesCurrencies_NotEnoughCurrencies_Error(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	ctx, r := gin.CreateTestContext(responseRecorder)
	repositoryMock := mocks.NewMockRatesRepository(t)
	service := domain.NewRatesService(repositoryMock)

	ratesHandler := NewRatesHandler(service, zap.NewNop())
	req, _ := http.NewRequest("GET", "/rates?currencies=EUR", nil)
	r.GET("/rates", ratesHandler.GetRatesCurrencies)
	ctx.Request = req

	r.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	repositoryMock.AssertExpectations(t)
}

func TestRates_GetRatesCurrencies_InvalidCurrency_Error(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	ctx, r := gin.CreateTestContext(responseRecorder)
	repositoryMock := mocks.NewMockRatesRepository(t)
	service := domain.NewRatesService(repositoryMock)

	ratesHandler := NewRatesHandler(service, zap.NewNop())
	req, _ := http.NewRequest("GET", "/rates?currencies=EUR,US", nil)
	r.GET("/rates", ratesHandler.GetRatesCurrencies)
	ctx.Request = req

	r.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	repositoryMock.AssertExpectations(t)
}

func TestRates_GetRatesCurrencies_ApiFailed_Error(t *testing.T) {
	responseRecorder := httptest.NewRecorder()

	ctx, r := gin.CreateTestContext(responseRecorder)
	repositoryMock := mocks.NewMockRatesRepository(t)
	service := domain.NewRatesService(repositoryMock)

	eur, _ := domain.NewCurrency("EUR")
	gbp, _ := domain.NewCurrency("GBP")
	currencies := []*domain.Currency{eur, gbp}

	baseRates := make(map[string]float64, 2)
	baseRates["EUR"] = 1
	baseRates["GBP"] = 2

	baseCurrency, _ := domain.NewCurrency("USD")

	repositoryMock.On("GetRates", mock.Anything, *baseCurrency, currencies).Return(
		nil,
		errors.New("failure"),
	)

	ratesHandler := NewRatesHandler(service, zap.NewNop())
	req, _ := http.NewRequest("GET", "/rates?currencies=EUR,GBP", nil)
	r.GET("/rates", ratesHandler.GetRatesCurrencies)
	ctx.Request = req

	r.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	repositoryMock.AssertExpectations(t)
}
