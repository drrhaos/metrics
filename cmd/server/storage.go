package main

import (
	"fmt"
	"strconv"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (storage *MemStorage) init() {
	storage.counter = make(map[string]int64)
	storage.gauge = make(map[string]float64)
}
func (storage *MemStorage) updateMetric(typeMetric string, nameMetric string, valueMetric string) bool {
	fmt.Println(typeMetric, nameMetric, valueMetric)
	if typeMetric == typeMetricCounter {
		valueIntMetric, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			return false
		}
		storage.counter[nameMetric] += valueIntMetric
	}

	if typeMetric == typeMetricGauge {
		valueFloatMetric, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			return false
		}
		storage.gauge[nameMetric] = valueFloatMetric
	}
	return true
}

func (storage *MemStorage) getMetric(typeMetric string, nameMetric string) (currentValue string, exists bool) {
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
