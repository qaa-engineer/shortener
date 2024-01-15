package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	Address         string
	BaseResponseURL string
}

type EnvConfiguration struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Address:         "localhost:8080",
		BaseResponseURL: "http://localhost:8080/",
	}
}

func (c *Configuration) ParseConfiguration() {

	cfg := EnvConfiguration{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Print(err)
	}

	if len(cfg.ServerAddress) > 0 {
		c.Address = cfg.ServerAddress
	}

	if len(cfg.BaseURL) > 0 {
		c.BaseResponseURL = cfg.BaseURL
	}

}
