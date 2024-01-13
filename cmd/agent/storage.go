package main

import (
	"log"
	"reflect"
	"runtime"
)

const randomValueName = "RandomValue"
const pollCountName = "PollCount"

var nameGauges = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (stat *MemStorage) updateGauge(nameMetric string, valueMetric float64) {
	stat.gauge[nameMetric] = valueMetric
}
func (stat *MemStorage) updateCounter(nameMetric string, valueMetric int64) {
	stat.counter[nameMetric] = valueMetric
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
		log.Println("Тип значения не соответствует uint", name)
		return floatValue, false
	}
	return floatValue, true
}
