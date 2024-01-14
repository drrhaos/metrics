package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

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
	for nameMetric, valueMetric := range metricsCPU.gauge {
		urlStr := fmt.Sprintf(urlUpdateGaugeConst, endpoint, nameMetric, valueMetric)
		r, _ := http.NewRequest(http.MethodPost, urlStr, nil)
		r.Header.Add("Content-Type", "text/plain")
		resp, err := client.Do(r)
		if err == nil {
			defer resp.Body.Close()
		} else {
			log.Println("Ошибка при выполнении запроса", urlStr)
		}
	}

	for nameMetric, valueMetric := range metricsCPU.counter {
		urlStr := fmt.Sprintf(urlUpdateCounterConst, endpoint, nameMetric, valueMetric)
		r, _ := http.NewRequest(http.MethodPost, urlStr, nil)
		r.Header.Add("Content-Type", "text/plain")
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
