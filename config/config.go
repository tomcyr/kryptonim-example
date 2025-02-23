package config

import (
	"fmt"
	"github.com/aranw/yamlcfg"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	OpenExchangeRates OpenExchangeRates `yaml:"open_exchange_rates" validate:"required"`
	REST              REST              `yaml:"rest" validate:"required"`
}

type OpenExchangeRates struct {
	AppID   string `yaml:"app_id" validate:"required"`
	BaseURL string `yaml:"base_url" validate:"required"`
}

type REST struct {
	Port int `yaml:"port" validate:"required"`
}

func (c Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

func ParseConfig(path string) (*Config, error) {
	cfg, err := yamlcfg.Parse[Config](path)
	if err != nil {
		return nil, fmt.Errorf("parsing yaml config failed: %w", err)
	}

	return cfg, nil
}
