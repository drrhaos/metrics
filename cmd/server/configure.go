package main

import (
	"flag"
	"net/url"

	"github.com/caarlos0/env/v10"
	"github.com/drrhaos/metrics/internal/logger"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int64  `env:"STORE_INTERVAL" envDefault:"-1"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE" envDefault:"true"`
}

func (cfg *Config) readStartParams() bool {
	err := env.Parse(cfg)
	if err != nil {
		logger.Log.Info("Не удалось найти переменные окружения ")
	}
	address := flag.String("a", "127.0.0.1:8080", "Сетевой адрес host:port")
	storeInterval := flag.Int64("i", 300, "Интервал сохранения показаний")
	fileStoragePath := flag.String("f", "/tmp/metrics-db.json", "Путь к файлу с показаниями")
	restore := flag.Bool("r", true, "Загружать последние сохранения показаний")
	flag.Parse()
	if cfg.Address == "" {
		cfg.Address = *address
	}

	if cfg.StoreInterval == -1 {
		cfg.StoreInterval = *storeInterval
	}

	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *fileStoragePath
	}

	if !cfg.Restore {
		cfg.Restore = *restore
	}

	_, errURL := url.ParseRequestURI("http://" + cfg.Address)
	if errURL != nil {
		flag.PrintDefaults()
		return false
	} else {
		return true
	}
}
