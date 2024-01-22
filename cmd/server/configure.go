package main

import (
	"flag"
	"net/url"

	"github.com/caarlos0/env/v6"
	"github.com/drrhaos/metrics/internal/logger"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func readStartParams() (Config, bool) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		logger.Log.Info("Не удалось найти переменные окружения")
	}
	address := flag.String("a", "127.0.0.1:8080", "Net address endpoint host:port")
	flag.Parse()
	if cfg.Address == "" {
		cfg.Address = *address
	}
	_, errURL := url.ParseRequestURI("http://" + cfg.Address)
	if errURL != nil {
		flag.PrintDefaults()
		return cfg, false
	} else {
		return cfg, true
	}
}
