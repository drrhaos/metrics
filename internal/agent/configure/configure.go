package configure

import (
	"encoding/json"
	"flag"
	"net/url"
	"os"

	"metrics/internal/logger"

	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

type Config struct {
	Address        string `env:"ADDRESS" json:"address,omitempty"`                                 // адрес эндпоинта HTTP-сервера
	CryptoKeyPath  string `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`                           // путь до файла с публичным ключом
	ConfigPath     string `env:"CONFIG"`                                                           // путь до файла с конфигурацией
	Key            string `env:"KEY" json:"key,omitempty"`                                         // ключ для проверки целостности данных в запросе
	PollInterval   int64  `env:"POLL_INTERVAL" envDefault:"-1" json:"poll_interval,omitempty"`     // частота опроса метрик из пакета runtime
	RateLimit      int    `env:"RATE_LIMIT" envDefault:"-1" json:"rate_limit,omitempty"`           // количество одновременно исходящих запросов на сервер ограничение «сверху»
	ReportInterval int64  `env:"REPORT_INTERVAL" envDefault:"-1" json:"report_interval,omitempty"` // частота отправки метрик на сервер
}

func (cfg *Config) readFlags() {
	address := flag.String("a", "", "Адрес эндпоинта HTTP-сервера host:port")
	reportInterval := flag.Int64("r", -1, "Частота отправки метрик на сервер")
	pollInterval := flag.Int64("p", -1, "Частота опроса метрик из пакета runtime")
	rateLimit := flag.Int("l", -1, "Количество одновременно исходящих запросов на сервер ограничение «сверху»")
	key := flag.String("k", "", "Ключ для проверки целостности данных в запросе")
	cryptoKeyPath := flag.String("crypto-key", "", "Путь до файла с публичным ключом")
	var configPath string
	flag.StringVar(&configPath, "config", "", "Путь до файла конфигурации")
	flag.StringVar(&configPath, "c", "", "Путь до файла конфигурации")
	flag.Parse()

	cfg.ConfigPath = configPath
	cfg.Address = *address
	cfg.ReportInterval = *reportInterval
	cfg.PollInterval = *pollInterval
	cfg.RateLimit = *rateLimit
	cfg.Key = *key
	cfg.CryptoKeyPath = *cryptoKeyPath
}

func (cfg *Config) readEnv() {
	var tmpConf Config
	err := env.Parse(&tmpConf)
	if err != nil {
		logger.Log.Warn("Не удалось найти переменные окружения", zap.Error(err))
	}

	if cfg.Address == "" {
		cfg.Address = tmpConf.Address
	}

	if cfg.Key == "" {
		cfg.Key = tmpConf.Key
	}

	if cfg.CryptoKeyPath == "" {
		cfg.CryptoKeyPath = tmpConf.CryptoKeyPath
	}

	if cfg.PollInterval == -1 {
		cfg.PollInterval = tmpConf.PollInterval
	}

	if cfg.ReportInterval == -1 {
		cfg.ReportInterval = tmpConf.ReportInterval
	}

	if cfg.RateLimit == -1 {
		cfg.RateLimit = tmpConf.RateLimit
	}
}

func (cfg *Config) readJSON() error {
	if cfg.ConfigPath == "" {
		return nil
	}
	dataJson, err := os.ReadFile(cfg.ConfigPath)
	if err != nil {
		return err
	}
	var tmpConf Config
	err = json.Unmarshal(dataJson, &tmpConf)
	if err != nil {
		return err
	}

	if cfg.Address == "" {
		cfg.Address = tmpConf.Address
	}

	if cfg.Key == "" {
		cfg.Key = tmpConf.Key
	}

	if cfg.CryptoKeyPath == "" {
		cfg.CryptoKeyPath = tmpConf.CryptoKeyPath
	}

	if cfg.PollInterval == -1 {
		cfg.PollInterval = tmpConf.PollInterval
	}

	if cfg.ReportInterval == -1 {
		cfg.ReportInterval = tmpConf.ReportInterval
	}

	if cfg.RateLimit == -1 {
		cfg.RateLimit = tmpConf.RateLimit
	}

	return nil
}

func (cfg *Config) checkConfig() (ok bool) {
	if cfg.Address == "" {
		cfg.Address = "127.0.0.1:8080"
	}

	if cfg.PollInterval == -1 {
		cfg.PollInterval = 2
	}

	if cfg.ReportInterval == -1 {
		cfg.ReportInterval = 10
	}

	if cfg.RateLimit == -1 {
		cfg.RateLimit = 1
	}

	_, errURL := url.ParseRequestURI("http://" + cfg.Address)
	if errURL != nil {
		logger.Log.Error("неверный формат адреса")
		return false
	}

	return true
}

func (cfg *Config) ReadConfig() bool {
	cfg.readFlags()
	cfg.readEnv()
	cfg.readJSON()
	cfg.checkConfig()
	return true
}
