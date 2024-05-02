package configure

import (
	"flag"
	"net/url"

	"metrics/internal/logger"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
	Key            string `env:"KEY"`
}

func (cfg *Config) ReadConfig() bool {
	err := env.Parse(cfg)
	if err != nil {
		logger.Log.Info("Не удалось найти переменные окружения")
	}

	address := flag.String("a", "127.0.0.1:8080", "Net address endpoint host:port")
	reportInterval := flag.Int64("r", 10, "Report interval integer sec > 0")
	pollInterval := flag.Int64("p", 2, "Pool interval integer sec > 0")
	rateLimit := flag.Int("l", 1, "number of simultaneous outgoing requests to the server")
	key := flag.String("k", "", "sha256 key for encryption of transmitted data")
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

	if cfg.RateLimit == 0 {
		cfg.RateLimit = *rateLimit
	}

	if cfg.Key == "" {
		cfg.Key = *key
	}
	_, errURL := url.ParseRequestURI("http://" + cfg.Address)

	if cfg.PollInterval <= 0 || cfg.ReportInterval <= 0 || errURL != nil {
		flag.PrintDefaults()
		return false
	} else {
		return true
	}
}
