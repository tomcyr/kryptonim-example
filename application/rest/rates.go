package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tomcyr/kryptonim-example/domain"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type RatesHandler struct {
	svc    domain.RatesService
	logger *zap.Logger
}

func NewRatesHandler(svc domain.RatesService, logger *zap.Logger) *RatesHandler {
	return &RatesHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h RatesHandler) GetRatesCurrencies(ctx *gin.Context) {
	symbols := strings.Split(ctx.Query("currencies"), ",")
	if len(symbols) < 2 {
		h.logger.Debug("not enough currencies")
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	currencies := make([]*domain.Currency, len(symbols))
	for k, symbol := range symbols {
		cur, err := domain.NewCurrency(symbol)
		if err != nil {
			h.logger.Debug("invalid currencies", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		currencies[k] = cur
	}
	rates, err := h.svc.GetRates(ctx.Request.Context(), currencies)
	if err != nil {
		h.logger.Error("an error occurred", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, rates)
}
