package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func updateMertics(metricsCPU *MemStorage, PollCount int64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metricsCPU.updateGauge(randomValueName, rand.Float64())
	metricsCPU.updateCounter(pollCountName, PollCount)

	for _, name := range nameGauges {
		floatValue, ok := getFloat64MemStats(m, name)
		if ok {
			metricsCPU.updateGauge(name, floatValue)
		}
	}
}

func sendMetrics(endpoint string, metricsCPU *MemStorage) {
	client := &http.Client{}
	currentGauges, ok := metricsCPU.getGauges()
	if !ok {
		return
	}
	for nameMetric, valueMetric := range currentGauges {
		var metric Metrics
		metric.MType = typeMetricGauge
		metric.ID = nameMetric
		metric.Value = &valueMetric
		urlStr := fmt.Sprintf(urlUpdateJSONConst, endpoint)
		reqData, err := json.Marshal(metric)
		if err != nil {
			log.Println(err)
			return
		}

		buf, err := compressReqData(reqData)
		if err != nil {
			log.Println(err)
			return
		}
		r, _ := http.NewRequest(http.MethodPost, urlStr, buf)
		r.Header.Set("Content-Type", "application/json")
		// r.Header.Set("Content-Encoding", "gzip")
		resp, err := client.Do(r)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

	}

	currentCounters, ok := metricsCPU.getCounters()
	if !ok {
		return
	}
	for nameMetric, valueMetric := range currentCounters {
		var metric Metrics
		metric.MType = typeMetricCounter
		metric.ID = nameMetric
		metric.Delta = &valueMetric
		urlStr := fmt.Sprintf(urlUpdateJSONConst, endpoint)
		reqData, err := json.Marshal(metric)
		if err != nil {
			log.Println("Ошибка при выполнении запроса", urlStr)
			return
		}
		buf, err := compressReqData(reqData)
		if err != nil {
			log.Println(err)
			return
		}
		r, _ := http.NewRequest(http.MethodPost, urlStr, buf)
		r.Header.Set("Content-Type", "application/json")
		// r.Header.Set("Content-Encoding", "gzip")
		resp, err := client.Do(r)
		if err == nil {
			defer resp.Body.Close()
		} else {
			log.Println("Ошибка при выполнении запроса", urlStr)
		}
	}
}

func collectMetrics(cfg Config) {
	metricsCPU := &MemStorage{
		counter: make(map[string]int64),
		gauge:   make(map[string]float64),
		mut:     sync.Mutex{},
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
			sendMetrics(cfg.Address, metricsCPU)
		}
	}()
	select {}
}
