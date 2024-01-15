package config

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	Address         string
	BaseResponseURL string
	FileStoragePath string
}

type EnvConfiguration struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfiguration() *Configuration {
	return &Configuration{
		Address:         "localhost:8080",
		BaseResponseURL: "http://localhost:8080/",
		FileStoragePath: "/tmp/short_url.json",
	}
}

func (c *Configuration) ParseConfiguration() {

	flag.Func("a", "адрес запуска HTTP-сервера (переменная SERVER_ADDRESS)", func(s string) error {
		_, err := url.ParseRequestURI(s)
		if err != nil {
			return err
		}

		c.Address = s

		return nil
	})

	flag.Func("b", "базовый адрес результирующего сокращённого URL (переменная BASE_URL)", func(s string) error {
		_, err := url.ParseRequestURI(s)
		if err != nil {
			return err
		}

		c.BaseResponseURL = s

		return nil
	})

	flag.Func("f", "путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH)", func(s string) error {

		if len(s) == 0 {
			c.FileStoragePath = ""
			return nil
		}

		c.FileStoragePath = s
		return nil
	})

	flag.Parse()

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

	if len(cfg.FileStoragePath) > 0 {
		c.FileStoragePath = cfg.FileStoragePath
	}

}
