package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (stat *MemStorage) update(cur runtime.MemStats) {
	stat.gauge["Alloc"] = float64(cur.Alloc)
	stat.gauge["BuckHashSys"] = float64(cur.BuckHashSys)
	stat.gauge["Frees"] = float64(cur.Frees)
	stat.gauge["GCCPUFraction"] = float64(cur.GCCPUFraction)
	stat.gauge["GCSys"] = float64(cur.GCSys)
	stat.gauge["HeapAlloc"] = float64(cur.HeapAlloc)
	stat.gauge["HeapIdle"] = float64(cur.HeapIdle)
	stat.gauge["HeapInuse"] = float64(cur.HeapInuse)
	stat.gauge["HeapObjects"] = float64(cur.HeapObjects)
	stat.gauge["HeapReleased"] = float64(cur.HeapReleased)
	stat.gauge["HeapSys"] = float64(cur.HeapSys)
	stat.gauge["LastGC"] = float64(cur.LastGC)
	stat.gauge["Lookups"] = float64(cur.Lookups)
	stat.gauge["MCacheInuse"] = float64(cur.MCacheInuse)
	stat.gauge["MCacheSys"] = float64(cur.MCacheSys)
	stat.gauge["MSpanInuse"] = float64(cur.MSpanInuse)
	stat.gauge["MSpanSys"] = float64(cur.MSpanSys)
	stat.gauge["Mallocs"] = float64(cur.Mallocs)
	stat.gauge["NextGC"] = float64(cur.NextGC)
	stat.gauge["NumForcedGC"] = float64(cur.NumForcedGC)
	stat.gauge["NumGC"] = float64(cur.NumGC)
	stat.gauge["OtherSys"] = float64(cur.OtherSys)
	stat.gauge["PauseTotalNs"] = float64(cur.PauseTotalNs)
	stat.gauge["StackInuse"] = float64(cur.StackInuse)
	stat.gauge["StackSys"] = float64(cur.StackSys)
	stat.gauge["Sys"] = float64(cur.Sys)
	stat.gauge["TotalAlloc"] = float64(cur.TotalAlloc)
}

func sendData(endpoint *string, typeMetric string, nameMetric string, valueMetric string) {
	client := &http.Client{}
	urlStr := fmt.Sprintf("http://%s/update/%s/%s/%s", *endpoint, typeMetric, nameMetric, valueMetric)
	r, _ := http.NewRequest(http.MethodPost, urlStr, nil) // URL-encoded payload
	r.Header.Add("Content-Type", "text/plain")
	_, err := client.Do(r)
	if err != nil {
		// fmt.Println("error")
	}
}

func main() {
	endpoint := flag.String("a", "127.0.0.1:8080", "Net address endpoint host:port")
	reportInterval := flag.Int64("r", 10, "Report interval")
	pollInterval := flag.Int64("p", 2, "Pool interval")
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err == nil {
		endpoint = &cfg.Address
		reportInterval = &cfg.ReportInterval
		pollInterval = &cfg.PollInterval
	}

	metricsCpu := MemStorage{
		gauge: map[string]float64{},
	}
	var PollCount int64 = 0
	var m runtime.MemStats
	var RandomValue float64
	for {
		PollCount++
		runtime.ReadMemStats(&m)
		metricsCpu.update(m)
		RandomValue = rand.Float64()

		if (PollCount**pollInterval)%*reportInterval == 0 {
			// update
			sendData(endpoint, "counter", "PollCount", strconv.FormatInt(PollCount, 10))
			sendData(endpoint, "gauge", "RandomValue", strconv.FormatFloat(RandomValue, 'f', -1, 64))
			for nameMetric, valueMetric := range metricsCpu.gauge {
				sendData(endpoint, "gauge", nameMetric, strconv.FormatFloat(valueMetric, 'f', -1, 64))
				// fmt.Println(named, value)
			}

			// fmt.Printf("alloc %f", mm.Alloc)
		}

		time.Sleep(time.Duration(*pollInterval) * time.Second)
	}
}
