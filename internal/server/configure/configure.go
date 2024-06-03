// Package configure предназначен для настройки программы.
package configure

import (
	"encoding/json"
	"flag"
	"net/url"
	"os"

	"metrics/internal/logger"

	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

// Config хранит текущую конфигурацию сервиса.
type Config struct {
	Address         string `env:"ADDRESS" json:"address,omitempty"`                               // адрес сервера
	ConfigPath      string `env:"CONFIG"`                                                         // путь до файла с конфигурацией
	CryptoKeyPath   string `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`                         // передайте путь до файла с приватным ключом
	DatabaseDsn     string `env:"DATABASE_DSN" json:"database_dsn,omitempty"`                     // DSN базы данных
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"nil" json:"store_file,omitempty"` // путь до хранилища данных
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet,omitempty"`                 // строковое представление бесклассовой адресации (CIDR)
	Key             string `env:"KEY" json:"key,omitempty"`                                       // ключ для проверки целостности данных в запросе
	Restore         bool   `env:"RESTORE" envDefault:"true" json:"restore,omitempty"`             // флаг указывающий на необходимость восстановления из хранилища при старте
	StoreInterval   int64  `env:"STORE_INTERVAL" envDefault:"-1" json:"store_interval,omitempty"` // интервал сохранения данных
}

func (cfg *Config) readFlags() {
	address := flag.String("a", "", "Сетевой адрес host:port")
	var configPath string
	flag.StringVar(&configPath, "config", "", "Путь до файла конфигурации")
	flag.StringVar(&configPath, "c", "", "Путь до файла конфигурации")
	cryptoKeyPath := flag.String("crypto-key", "", "Путь до файла с приватным ключом")
	databaseDsn := flag.String("d", "",
		"Сетевой адрес базя данных postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable")
	fileStoragePath := flag.String("f", "nil", "Путь к файлу с показаниями")
	key := flag.String("k", "", "ключ для проверки целостности данных в запросе")
	restore := flag.Bool("r", true, "Загружать последние сохранения показаний")
	storeInterval := flag.Int64("i", -1, "Интервал сохранения показаний")
	trustedSubnet := flag.String("t", "", "Строковое представление бесклассовой адресации (CIDR)")
	flag.Parse()

	cfg.Address = *address

	cfg.ConfigPath = configPath
	cfg.CryptoKeyPath = *cryptoKeyPath
	cfg.DatabaseDsn = *databaseDsn
	cfg.FileStoragePath = *fileStoragePath
	cfg.Key = *key
	if *fileStoragePath == "" {
		cfg.Restore = false
	} else {
		cfg.Restore = *restore
	}
	cfg.StoreInterval = *storeInterval
	cfg.TrustedSubnet = *trustedSubnet
}

func (cfg *Config) readEnv() error {
	var tmpCfg Config
	err := env.Parse(&tmpCfg)
	if err != nil {
		logger.Log.Warn("Не удалось найти переменные окружения", zap.Error(err))
		return err
	}

	if tmpCfg.Address != "" {
		cfg.Address = tmpCfg.Address
	}
	if tmpCfg.ConfigPath != "" {
		cfg.ConfigPath = tmpCfg.ConfigPath
	}
	if tmpCfg.CryptoKeyPath != "" {
		cfg.CryptoKeyPath = tmpCfg.CryptoKeyPath
	}
	if tmpCfg.DatabaseDsn != "" {
		cfg.DatabaseDsn = tmpCfg.DatabaseDsn
	}
	if tmpCfg.FileStoragePath != "nil" {
		cfg.FileStoragePath = tmpCfg.FileStoragePath
	}
	if tmpCfg.Key != "" {
		cfg.Key = tmpCfg.Key
	}

	cfg.Restore = tmpCfg.Restore

	if tmpCfg.StoreInterval != -1 {
		cfg.StoreInterval = tmpCfg.StoreInterval
	}

	if tmpCfg.TrustedSubnet != "" {
		cfg.TrustedSubnet = tmpCfg.TrustedSubnet
	}

	return nil
}

func (cfg *Config) readJSON() error {
	if cfg.ConfigPath == "" {
		return nil
	}

	dataJSON, err := os.ReadFile(cfg.ConfigPath)
	if err != nil {
		logger.Log.Error("не удалось прочитать файл конфигурации", zap.Error(err))
		return err
	}
	var tmpCfg Config
	err = json.Unmarshal(dataJSON, &tmpCfg)
	if err != nil {
		logger.Log.Error("файл конфигурации имеет не верный формат", zap.Error(err))
		return err
	}

	if tmpCfg.Address != "" && cfg.Address == "" {
		cfg.Address = tmpCfg.Address
	}
	if tmpCfg.ConfigPath != "" && cfg.ConfigPath == "" {
		cfg.ConfigPath = tmpCfg.ConfigPath
	}
	if tmpCfg.CryptoKeyPath != "" && cfg.CryptoKeyPath == "" {
		cfg.CryptoKeyPath = tmpCfg.CryptoKeyPath
	}
	if tmpCfg.DatabaseDsn != "" && cfg.DatabaseDsn == "" {
		cfg.DatabaseDsn = tmpCfg.DatabaseDsn
	}
	if tmpCfg.FileStoragePath != "nil" && cfg.FileStoragePath == "nil" {
		cfg.FileStoragePath = tmpCfg.FileStoragePath
	}
	if tmpCfg.Key != "" && cfg.Key == "" {
		cfg.Key = tmpCfg.Key
	}

	cfg.Restore = tmpCfg.Restore

	if tmpCfg.StoreInterval != -1 && cfg.StoreInterval == -1 {
		cfg.StoreInterval = tmpCfg.StoreInterval
	}

	if tmpCfg.TrustedSubnet != "" && cfg.TrustedSubnet == "" {
		cfg.TrustedSubnet = tmpCfg.TrustedSubnet
	}

	return nil
}

func (cfg *Config) checkConfig() bool {
	if cfg.Address == "" {
		cfg.Address = "127.0.0.1:8080"
	}

	if cfg.StoreInterval == -1 {
		cfg.StoreInterval = 300
	}

	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = "/tmp/metrics-db.json"
	}

	_, errURL := url.ParseRequestURI("http://" + cfg.Address)
	if errURL != nil {
		logger.Log.Error("неверный формат адреса")
		return false
	}

	return true
}

// ReadConfig читает конфигурацию сервера
func (cfg *Config) ReadConfig() bool {
	cfg.readFlags()

	err := cfg.readEnv()
	if err != nil {
		return false
	}

	err = cfg.readJSON()
	if err != nil {
		return false
	}
	cfg.checkConfig()

	return true
}
