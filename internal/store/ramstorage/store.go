// Package ramstorage реализует хранение метрик в памяти.
package ramstorage

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"metrics/internal/logger"
	"metrics/internal/store"
)

// RAMStorage хранилище метрик.
type RAMStorage struct {
	Gauge   map[string]float64 `json:"gauge"`   // набор метрик counter
	Counter map[string]int64   `json:"counter"` // набор метрик gauge
	Mut     sync.Mutex         // мютекс
}

// NewStorage инициализарует хранилище.
func NewStorage() *RAMStorage {
	return &RAMStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
		Mut:     sync.Mutex{},
	}
}

// SaveMetrics сохраняет метрики.
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
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o666)
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

// LoadMetrics загружает метрики.
func (storage *RAMStorage) LoadMetrics(filePath string) bool {
	if storage == nil {
		return false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0o666)
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

// UpdateCounter обновляет метрику counter.
func (storage *RAMStorage) UpdateCounter(_ context.Context, nameMetric string, valueMetric int64) bool {
	if storage == nil {
		return false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	storage.Counter[nameMetric] += valueMetric
	return true
}

// UpdateGauge обновляет метрику gauge.
func (storage *RAMStorage) UpdateGauge(_ context.Context, nameMetric string, valueMetric float64) bool {
	if storage == nil {
		return false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	storage.Gauge[nameMetric] = valueMetric
	return true
}

// GetGauges возвращает метрики gauges.
func (storage *RAMStorage) GetGauges(_ context.Context) (map[string]float64, bool) {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	return storage.Gauge, true
}

// GetCounters возвращает метрики counters.
func (storage *RAMStorage) GetCounters(_ context.Context) (map[string]int64, bool) {
	if storage == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	return storage.Counter, true
}

// GetCounter возвращает метрику counters.
func (storage *RAMStorage) GetCounter(_ context.Context, nameMetric string) (currentValue int64, exists bool) {
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

// GetGauge возвращает метрику gauge.
func (storage *RAMStorage) GetGauge(_ context.Context, nameMetric string) (currentValue float64, exists bool) {
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

// GetBatchMetrics возвращает пачку хранимых метрик.
func (storage *RAMStorage) GetBatchMetrics(_ context.Context) (metrics []store.Metrics, exists bool) {
	storage.Mut.Lock()
	defer storage.Mut.Unlock()

	for nameMetric, valueMetric := range storage.Gauge {
		hd := valueMetric
		metrics = append(metrics, store.Metrics{ID: nameMetric, MType: "gauge", Value: &hd})
	}

	for nameMetric, valueMetric := range storage.Counter {
		hd := valueMetric
		metrics = append(metrics, store.Metrics{ID: nameMetric, MType: "counter", Delta: &hd})
	}
	return metrics, true
}

// Ping проверяет доступность хранилища.
func (storage *RAMStorage) Ping(_ context.Context) bool {
	return false
}
