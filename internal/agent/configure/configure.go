package configure

import (
	"flag"
	"net/url"

	"metrics/internal/logger"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS"`         // адрес эндпоинта HTTP-сервера
	CryptoKeyPath  string `env:"CRYPTO_KEY"`      // путь до файла с публичным ключом
	Key            string `env:"KEY"`             // ключ для проверки целостности данных в запросе
	PollInterval   int64  `env:"POLL_INTERVAL"`   // частота опроса метрик из пакета runtime
	RateLimit      int    `env:"RATE_LIMIT"`      // количество одновременно исходящих запросов на сервер ограничение «сверху»
	ReportInterval int64  `env:"REPORT_INTERVAL"` // частота отправки метрик на сервер
}

func (cfg *Config) ReadConfig() bool {
	err := env.Parse(cfg)
	if err != nil {
		logger.Log.Info("Не удалось найти переменные окружения")
	}

	address := flag.String("a", "127.0.0.1:8080", "Адрес эндпоинта HTTP-сервера host:port")
	reportInterval := flag.Int64("r", 10, "Частота отправки метрик на сервер")
	pollInterval := flag.Int64("p", 2, "Частота опроса метрик из пакета runtime")
	rateLimit := flag.Int("l", 1, "Количество одновременно исходящих запросов на сервер ограничение «сверху»")
	key := flag.String("k", "", "Ключ для проверки целостности данных в запросе")
	cryptoKeyPath := flag.String("crypto-key", "", "Путь до файла с публичным ключом")
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

	if cfg.CryptoKeyPath == "" {
		cfg.CryptoKeyPath = *cryptoKeyPath
	}

	_, errURL := url.ParseRequestURI("http://" + cfg.Address)

	if cfg.PollInterval <= 0 || cfg.ReportInterval <= 0 || errURL != nil {
		flag.PrintDefaults()
		return false
	}
	return true
}
