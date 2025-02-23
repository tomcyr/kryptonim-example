package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tomcyr/kryptonim-example/application/rest"
	"github.com/tomcyr/kryptonim-example/config"
	"github.com/tomcyr/kryptonim-example/domain"
	"github.com/tomcyr/kryptonim-example/infrastructure"
	"go.uber.org/zap"
)

var configFile = flag.String(
	"config_file",
	"config/config.yaml",
	"Path to the YAML config",
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	flag.Parse()

	cfg, err := config.ParseConfig(*configFile)
	if err != nil {
		panic(err)
	}

	oxrRepo := infrastructure.NewLoggingRatesRepository(
		infrastructure.NewOpenExchangeRatesRepository(cfg.OpenExchangeRates.AppID, cfg.OpenExchangeRates.BaseURL),
		logger,
	)
	service := domain.NewRatesService(oxrRepo)

	router := gin.Default()
	ratesHandler := rest.NewRatesHandler(service, logger)
	router.GET("/rates", ratesHandler.GetRatesCurrencies)

	err = router.Run(fmt.Sprintf(":%d", cfg.REST.Port))
	if err != nil {
		panic(err)
	}
}
