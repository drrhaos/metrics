package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/drrhaos/metrics/internal/logger"
)

type MemStorage struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
	Mut     sync.Mutex
}

func (storage *MemStorage) SaveMetrics(filePath string) bool {
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

func (storage *MemStorage) LoadMetrics(filePath string) bool {
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

func (storage *MemStorage) UpdateCounter(nameMetric string, valueMetric int64) bool {
	if storage == nil {
		return false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	storage.Counter[nameMetric] += valueMetric
	return true
}

func (storage *MemStorage) UpdateGauge(nameMetric string, valueMetric float64) bool {
	if storage == nil {
		return false
	}
	storage.Mut.Lock()
	defer storage.Mut.Unlock()
	storage.Gauge[nameMetric] = valueMetric
	return true
}

func (stat *MemStorage) GetGauges() (map[string]float64, bool) {
	if stat == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}
	stat.Mut.Lock()
	defer stat.Mut.Unlock()
	return stat.Gauge, true
}

func (stat *MemStorage) GetCounters() (map[string]int64, bool) {
	if stat == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}
	stat.Mut.Lock()
	defer stat.Mut.Unlock()
	return stat.Counter, true
}

func (storage *MemStorage) GetCounter(nameMetric string) (currentValue int64, exists bool) {
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

func (storage *MemStorage) GetGauge(nameMetric string) (currentValue float64, exists bool) {
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
