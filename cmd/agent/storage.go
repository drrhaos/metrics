package main

import (
	"reflect"
	"runtime"
	"sync"

	"github.com/drrhaos/metrics/internal/logger"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mut     sync.Mutex
}

func (stat *MemStorage) updateGauge(nameMetric string, valueMetric float64) {
	if stat == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return
	}
	stat.mut.Lock()
	defer stat.mut.Unlock()
	stat.gauge[nameMetric] = valueMetric
}

func (stat *MemStorage) updateCounter(nameMetric string, valueMetric int64) {
	if stat == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return
	}
	stat.mut.Lock()
	defer stat.mut.Unlock()
	stat.counter[nameMetric] = valueMetric
}

func (stat *MemStorage) getGauges() (map[string]float64, bool) {
	if stat == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}
	stat.mut.Lock()
	defer stat.mut.Unlock()
	return stat.gauge, true
}

func (stat *MemStorage) getCounters() (map[string]int64, bool) {
	if stat == nil {
		logger.Log.Panic("Хранилище не может быть nil")
		return nil, false
	}
	stat.mut.Lock()
	defer stat.mut.Unlock()
	return stat.counter, true
}

func getFloat64MemStats(m runtime.MemStats, name string) (float64, bool) {
	value := reflect.ValueOf(m).FieldByName(name)
	var floatValue float64
	switch value.Kind() {
	case reflect.Uint64:
		floatValue = float64(value.Uint())
	case reflect.Uint32:
		floatValue = float64(value.Uint())
	case reflect.Float64:
		floatValue = value.Float()
	default:
		logger.Log.Info("Тип значения не соответствует uint")
		return floatValue, false
	}
	return floatValue, true
}
