package main

import (
	"github.com/drrhaos/metrics/internal/database"
	"github.com/drrhaos/metrics/internal/logger"
	"github.com/drrhaos/metrics/internal/storage"
	"go.uber.org/zap"
)

type RouterStorage struct {
	RAM     storage.MemStorage
	DB      database.Database
	usageDB bool
}

func NewRouterStorage() *RouterStorage {
	switchStorage := &RouterStorage{RAM: *storage.NewStorage()}
	if cfg.DatabaseDsn != "" {
		switchStorage.DB = *database.NewDatabase()
		err := switchStorage.DB.Connect(cfg.DatabaseDsn)
		if err != nil {
			logger.Log.Panic("Не удалось подключиться к БД", zap.Error(err))
		}

		err = switchStorage.DB.Migrations()
		if err != nil {
			logger.Log.Panic("Не создать таблицы", zap.Error(err))
		}
		switchStorage.usageDB = true
		logger.Log.Info("Соединение с базой успешно установлено")
	}
	return switchStorage
}

func (storage *RouterStorage) SaveMetrics(filePath string) bool {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return false
	}
	return storage.RAM.SaveMetrics(filePath)
}

func (storage *RouterStorage) LoadMetrics(filePath string) bool {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return false
	}
	return storage.RAM.LoadMetrics(filePath)
}

func (storage *RouterStorage) UpdateCounter(nameMetric string, valueMetric int64) bool {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return false
	}
	if storage.usageDB {
		storage.DB.UpdateCounter(nameMetric, valueMetric)
	}

	return storage.RAM.UpdateCounter(nameMetric, valueMetric)
}

func (storage *RouterStorage) UpdateGauge(nameMetric string, valueMetric float64) bool {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return false
	}
	if storage.usageDB {
		storage.DB.UpdateGauge(nameMetric, valueMetric)
	}

	return storage.RAM.UpdateGauge(nameMetric, valueMetric)
}

func (storage *RouterStorage) GetGauges() (map[string]float64, bool) {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}

	valuesMetric, ok := storage.RAM.GetGauges()

	if storage.usageDB {
		valuesMetric, ok = storage.DB.GetGauges()
	}

	return valuesMetric, ok
}

func (storage *RouterStorage) GetCounters() (map[string]int64, bool) {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}
	valuesMetric, ok := storage.RAM.GetCounters()

	if storage.usageDB {
		valuesMetric, ok = storage.DB.GetCounters()
	}

	return valuesMetric, ok
}

func (storage *RouterStorage) GetCounter(nameMetric string) (currentValue int64, exists bool) {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return currentValue, false
	}

	currentValue, exists = storage.RAM.GetCounter(nameMetric)

	if storage.usageDB {
		currentValue, exists = storage.DB.GetCounter(nameMetric)
	}

	return currentValue, exists
}

func (storage *RouterStorage) GetGauge(nameMetric string) (currentValue float64, exists bool) {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return currentValue, false
	}
	currentValue, exists = storage.RAM.GetGauge(nameMetric)

	if storage.usageDB {
		currentValue, exists = storage.DB.GetGauge(nameMetric)
	}

	return currentValue, exists
}
