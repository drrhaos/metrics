package ramstorage

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"metrics/internal/logger"
)

type RAMStorage struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
	Mut     sync.Mutex
}

func NewStorage() *RAMStorage {
	return &RAMStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
		Mut:     sync.Mutex{},
	}
}

func (storage *RAMStorage) SaveMetrics(filePath string) bool {
	if storage == nil {
		return false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	data, err := json.Marshal(storage)
	if err != nil {
		logger.Log.Warn("не удалось преобразовать структуру")
		return false
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Warn("не удалось открыть файл")
		return false
	}

	_, err = file.Write(data)
	if err != nil {
		logger.Log.Warn("не удалось записать данные")
		return false
	}
	return true
}

func (storage *RAMStorage) LoadMetrics(filePath string) bool {
	if storage == nil {
		return false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Warn("не удалось открыть файл")
		return false
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		logger.Log.Warn("не удалось получить информацию о файле")
		return false
	}

	data := make([]byte, stat.Size())

	_, err = file.Read(data)
	if err != nil {
		logger.Log.Warn("не удалось записать данные")
		return false
	}
	file.Close()
	if err := json.Unmarshal(data, storage); err != nil {
		logger.Log.Warn("не удалось преобразовать структуру")
		return false
	}
	return true
}

func (storage *RAMStorage) UpdateCounter(ctx context.Context, nameMetric string, valueMetric int64) bool {
	if storage == nil {
		return false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	storage.Counter[nameMetric] += valueMetric
	return true
}

func (storage *RAMStorage) UpdateGauge(ctx context.Context, nameMetric string, valueMetric float64) bool {
	if storage == nil {
		return false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	storage.Gauge[nameMetric] = valueMetric
	return true
}

func (storage *RAMStorage) GetGauges(ctx context.Context) (map[string]float64, bool) {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	return storage.Gauge, true
}

func (storage *RAMStorage) GetCounters(ctx context.Context) (map[string]int64, bool) {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	return storage.Counter, true
}

func (storage *RAMStorage) GetCounter(ctx context.Context, nameMetric string) (currentValue int64, exists bool) {
	if storage == nil {
		return currentValue, false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	currentValue, ok := storage.Counter[nameMetric]
	if ok {
		exists = true
	}

	return currentValue, exists
}

func (storage *RAMStorage) GetGauge(ctx context.Context, nameMetric string) (currentValue float64, exists bool) {
	if storage == nil {
		return currentValue, false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	currentValue, ok := storage.Gauge[nameMetric]
	if ok {
		exists = true
	}

	return currentValue, exists
}

func (storage *RAMStorage) Ping(ctx context.Context) bool {
	return false
}
