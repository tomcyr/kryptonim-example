package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tomcyr/kryptonim-example/application/dto"
	"github.com/tomcyr/kryptonim-example/domain"
	"net/http"
	"strconv"
)

type ExchangeHandler struct {
	exchangeService domain.ExchangeService
}

func NewExchangeHandler(exchangeService domain.ExchangeService) *ExchangeHandler {
	return &ExchangeHandler{
		exchangeService: exchangeService,
	}
}

func (h ExchangeHandler) GetExchange(ctx *gin.Context) {
	from := ctx.Query("from")
	to := ctx.Query("to")
	amountStr := ctx.Query("amount")
	if len(from) == 0 || len(to) == 0 || len(amountStr) == 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount == 0 {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	exchange, err := domain.NewExchange(from, to, amount)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	res, err := h.exchangeService.Exchange(ctx.Request.Context(), exchange)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, dto.NewExchangeResponse(res))
}
