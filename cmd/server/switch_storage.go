package main

import (
	"github.com/drrhaos/metrics/internal/database"
	"github.com/drrhaos/metrics/internal/logger"
	"github.com/drrhaos/metrics/internal/storage"
)

type SwitchStorage struct {
	RAM     storage.MemStorage
	DB      database.Database
	usageDB bool
}

func NewSwitchStorage() *SwitchStorage {
	switchStorage := &SwitchStorage{RAM: *storage.NewStorage()}
	if cfg.DatabaseDsn != "" {
		switchStorage.DB = *database.NewDatabase()
		err := switchStorage.DB.Connect(cfg.DatabaseDsn)
		if err == nil {
			switchStorage.usageDB = true
		}
	}
	return switchStorage
}

func (storage *SwitchStorage) SaveMetrics(filePath string) bool {
	if storage == nil {
		return false
	}
	return storage.RAM.SaveMetrics(filePath)
}

func (storage *SwitchStorage) LoadMetrics(filePath string) bool {
	if storage == nil {
		return false
	}
	return storage.RAM.LoadMetrics(filePath)
}

func (storage *SwitchStorage) UpdateCounter(nameMetric string, valueMetric int64) bool {
	if storage == nil {
		return false
	}

	return storage.RAM.UpdateCounter(nameMetric, valueMetric)
}

func (storage *SwitchStorage) UpdateGauge(nameMetric string, valueMetric float64) bool {
	if storage == nil {
		return false
	}

	return storage.RAM.UpdateGauge(nameMetric, valueMetric)
}

func (storage *SwitchStorage) GetGauges() (map[string]float64, bool) {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}

	return storage.RAM.GetGauges()
}

func (storage *SwitchStorage) GetCounters() (map[string]int64, bool) {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}

	return storage.RAM.GetCounters()
}

func (storage *SwitchStorage) GetCounter(nameMetric string) (currentValue int64, exists bool) {
	if storage == nil {
		return currentValue, false
	}

	return storage.RAM.GetCounter(nameMetric)
}

func (storage *SwitchStorage) GetGauge(nameMetric string) (currentValue float64, exists bool) {
	if storage == nil {
		return currentValue, false
	}

	return storage.RAM.GetGauge(nameMetric)
}
