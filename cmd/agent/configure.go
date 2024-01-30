package main

import (
	"flag"
	"net/url"

	"github.com/caarlos0/env/v6"
	"github.com/drrhaos/metrics/internal/logger"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
}

func (cfg *Config) readStartParams() bool {
	err := env.Parse(cfg)
	if err != nil {
		logger.Log.Info("Не удалось найти переменные окружения")
	}

	address := flag.String("a", "127.0.0.1:8080", "Net address endpoint host:port")
	reportInterval := flag.Int64("r", 10, "Report interval integer sec > 0")
	pollInterval := flag.Int64("p", 2, "Pool interval integer sec > 0")
	flag.Parse()
	if cfg.Address == "" {
		cfg.Address = *address
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = *reportInterval
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = *pollInterval
	}
	_, errURL := url.ParseRequestURI("http://" + cfg.Address)

	if cfg.PollInterval <= 0 || cfg.ReportInterval <= 0 || errURL != nil {
		flag.PrintDefaults()
		return false
	} else {
		return true
	}
}
