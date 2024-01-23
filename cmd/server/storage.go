package main

import (
	"strconv"
	"sync"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mut     sync.Mutex
}

func (storage *MemStorage) updateCounter(nameMetric string, valueMetric int64) bool {
	if storage == nil {
		return false
	}
	storage.mut.Lock()
	defer storage.mut.Unlock()
	storage.counter[nameMetric] += valueMetric
	return true
}

func (storage *MemStorage) updateGauge(nameMetric string, valueMetric float64) bool {
	if storage == nil {
		return false
	}
	storage.mut.Lock()
	defer storage.mut.Unlock()
	storage.gauge[nameMetric] = valueMetric
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
		cur, ok := storage.counter[nameMetric]
		if ok {
			currentValue = strconv.FormatInt(cur, 10)
			exists = true
		}
	case typeMetricGauge:
		cur, ok := storage.gauge[nameMetric]
		if ok {
			currentValue = strconv.FormatFloat(cur, 'f', -1, 64)
			exists = true
		}
	}

	return currentValue, exists
}
