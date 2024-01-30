package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/drrhaos/metrics/internal/logger"
)

type MemStorage struct {
	Gauge   map[string]float64 `json:"gauge"`
	Counter map[string]int64   `json:"counter"`
	mut     sync.Mutex
}

func (storage *MemStorage) saveMetrics(filePath string) bool {
	if storage == nil {
		return false
	}
	storage.mut.Lock()
	defer storage.mut.Unlock()

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

func (storage *MemStorage) loadMetrics(filePath string) bool {
	if storage == nil {
		return false
	}
	fmt.Println(filePath)
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

func (storage *MemStorage) updateCounter(nameMetric string, valueMetric int64) bool {
	if storage == nil {
		return false
	}
	storage.mut.Lock()
	defer storage.mut.Unlock()
	storage.Counter[nameMetric] += valueMetric
	return true
}

func (storage *MemStorage) updateGauge(nameMetric string, valueMetric float64) bool {
	if storage == nil {
		return false
	}
	storage.mut.Lock()
	defer storage.mut.Unlock()
	storage.Gauge[nameMetric] = valueMetric
	return true
}

func (storage *MemStorage) getMetric(typeMetric string, nameMetric string) (currentValue string, exists bool) {
	if storage == nil {
		return currentValue, false
	}
	storage.mut.Lock()
	defer storage.mut.Unlock()
	switch typeMetric {
	case typeMetricCounter:
		cur, ok := storage.Counter[nameMetric]
		if ok {
			currentValue = strconv.FormatInt(cur, 10)
			exists = true
		}
	case typeMetricGauge:
		cur, ok := storage.Gauge[nameMetric]
		if ok {
			currentValue = strconv.FormatFloat(cur, 'f', -1, 64)
			exists = true
		}
	}

	return currentValue, exists
}
