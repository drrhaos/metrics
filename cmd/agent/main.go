package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)

func sendMetrics(endpoint string, metricsCPU MemStorage) {
	client := &http.Client{}
	for nameMetric, valueMetric := range metricsCPU.gauge {
		urlStr := fmt.Sprintf("http://%s/update/gauge/%s/%f", endpoint, nameMetric, valueMetric)
		r, _ := http.NewRequest(http.MethodPost, urlStr, nil) // URL-encoded payload
		r.Header.Add("Content-Type", "text/plain")
		resp, err := client.Do(r)
		if err == nil {
			defer resp.Body.Close()
		} else {
			log.Println("Ошибка при выполнении запроса", urlStr)
		}
	}

	for nameMetric, valueMetric := range metricsCPU.counter {
		urlStr := fmt.Sprintf("http://%s/update/counter/%s/%d", endpoint, nameMetric, valueMetric)
		r, _ := http.NewRequest(http.MethodPost, urlStr, nil) // URL-encoded payload
		r.Header.Add("Content-Type", "text/plain")
		resp, err := client.Do(r)
		if err == nil {
			defer resp.Body.Close()
		} else {
			log.Println("Ошибка при выполнении запроса", urlStr)
		}
	}
}

func main() {
	cfg, ok := readStartParams()
	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}
	var metricsCPU MemStorage
	metricsCPU.init()
	var PollCount int64 = 0
	var m runtime.MemStats

	for {
		runtime.ReadMemStats(&m)
		PollCount++
		metricsCPU.updateGauge(randomValueName, rand.Float64())
		metricsCPU.updateCounter(pollCountName, PollCount)

		for _, name := range nameGauges {
			floatValue, ok := getFloat64MemStats(m, name)
			if ok {
				metricsCPU.updateGauge(name, floatValue)
			}
		}

		if (PollCount*cfg.PollInterval)%cfg.ReportInterval == 0 {
			sendMetrics(cfg.Address, metricsCPU)
		}

		time.Sleep(time.Duration(cfg.PollInterval) * time.Second)
	}
}
