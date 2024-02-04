package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/drrhaos/metrics/internal/logger"
	"github.com/drrhaos/metrics/internal/storage"
	"go.uber.org/zap"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
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

func updateMertics(metricsCPU *storage.MemStorage, PollCount int64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metricsCPU.UpdateGauge(randomValueName, rand.Float64())
	metricsCPU.UpdateCounter(pollCountName, PollCount)

	for _, name := range nameGauges {
		floatValue, ok := getFloat64MemStats(m, name)
		if ok {
			metricsCPU.UpdateGauge(name, floatValue)
		}
	}
}

func sendMetric(metric Metrics) {
	client := &http.Client{}

	urlStr := fmt.Sprintf(urlUpdateJSONConst, cfg.Address)
	reqData, err := json.Marshal(metric)
	if err != nil {
		logger.Log.Warn("Не удалось создать JSON", zap.Error(err))
		return
	}

	buf, err := compressReqData(reqData)
	if err != nil {
		logger.Log.Warn("Не сжать данные", zap.Error(err))
		return
	}
	r, _ := http.NewRequest(http.MethodPost, urlStr, buf)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Content-Encoding", "gzip")
	resp, err := client.Do(r)
	if err != nil {
		logger.Log.Warn("Не удалось отправить запрос", zap.Error(err))
		return
	}
	defer resp.Body.Close()
}

func sendMetrics(metricsCPU *storage.MemStorage) {
	currentGauges, ok := metricsCPU.GetGauges()
	if !ok {
		return
	}
	for nameMetric, valueMetric := range currentGauges {
		var metric Metrics
		metric.MType = typeMetricGauge
		metric.ID = nameMetric
		metric.Value = &valueMetric
		sendMetric(metric)
	}

	currentCounters, ok := metricsCPU.GetCounters()
	if !ok {
		return
	}
	for nameMetric, valueMetric := range currentCounters {
		var metric Metrics
		metric.MType = typeMetricCounter
		metric.ID = nameMetric
		metric.Delta = &valueMetric
		sendMetric(metric)
	}
}

func collectMetrics() {
	metricsCPU := &storage.MemStorage{
		Counter: make(map[string]int64),
		Gauge:   make(map[string]float64),
		Mut:     sync.Mutex{},
	}

	go func() {
		var PollCount int64
		for {
			PollCount++
			time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
			updateMertics(metricsCPU, PollCount)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(cfg.ReportInterval) * time.Second)
			sendMetrics(metricsCPU)
		}
	}()
	select {}
}
