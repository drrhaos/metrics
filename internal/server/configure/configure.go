// Модуль configure предназначен для настройки программы.
package configure

import (
	"flag"
	"net/url"

	"metrics/internal/logger"

	"github.com/caarlos0/env/v10"
)

// Config хранит текущую конфигурацию сервиса.
type Config struct {
	Address         string `env:"ADDRESS"`                        // адрес сервера
	StoreInterval   int64  `env:"STORE_INTERVAL" envDefault:"-1"` // интервал сохранения данных
	FileStoragePath string `env:"FILE_STORAGE_PATH"`              // путь до хранилища данных
	Restore         bool   `env:"RESTORE" envDefault:"true"`      // флаг указывающий на необходимость восстановления из хранилища при старте
	DatabaseDsn     string `env:"DATABASE_DSN"`                   // DSN базы данных
	Key             string `env:"KEY"`                            // ключ для проверки целостности данных в запросе
}

// ReadStartParams чтение настроек.
func (cfg *Config) ReadStartParams() bool {
	err := env.Parse(cfg)

	if err != nil {
		logger.Log.Info("Не удалось найти переменные окружения ")
	}

	flagConf := flag.NewFlagSet("address", flag.ExitOnError)
	address := flagConf.String("a", "127.0.0.1:8080", "Сетевой адрес host:port")
	storeInterval := flagConf.Int64("i", 300, "Интервал сохранения показаний")
	fileStoragePath := flagConf.String("f", "/tmp/metrics-db.json", "Путь к файлу с показаниями")
	restore := flagConf.Bool("r", true, "Загружать последние сохранения показаний")
	databaseDsn := flagConf.String("d", "",
		"Сетевой адрес базя данных postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable")
	key := flagConf.String("k", "", "sha256 key for encryption of transmitted data")

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
