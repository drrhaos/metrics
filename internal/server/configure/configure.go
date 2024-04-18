package configure

import (
	"flag"
	"net/url"

	"metrics/internal/logger"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int64  `env:"STORE_INTERVAL" envDefault:"-1"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE" envDefault:"true"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
}

func (cfg *Config) ReadStartParams() bool {
	err := env.Parse(cfg)
	if err != nil {
		logger.Log.Info("Не удалось найти переменные окружения ")
	}
	address := flag.String("a", "127.0.0.1:8080", "Сетевой адрес host:port")
	storeInterval := flag.Int64("i", 300, "Интервал сохранения показаний")
	fileStoragePath := flag.String("f", "/tmp/metrics-db.json", "Путь к файлу с показаниями")
	restore := flag.Bool("r", true, "Загружать последние сохранения показаний")
	databaseDsn := flag.String("d", "",
		"Сетевой адрес базя данных postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable")
	key := flag.String("k", "", "sha256 key for encryption of transmitted data")

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

	if cfg.DatabaseDsn == "" {
		cfg.DatabaseDsn = *databaseDsn
	}

	if cfg.Key == "" {
		cfg.Key = *key
	}

	_, errURL := url.ParseRequestURI("http://" + cfg.Address)
	if errURL != nil {
		flag.PrintDefaults()
		return false
	} else {
		return true
	}
}
