package main

import (
	"flag"
	"os"

	"github.com/drrhaos/metrics/internal/logger"
)

const typeMetricCounter = "counter"
const typeMetricGauge = "gauge"
const randomValueName = "RandomValue"
const pollCountName = "PollCount"

const urlUpdateJSONConst = "http://%s/update/"

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

const flagLogLevel = "info"

var cfg Config

func main() {
	if err := logger.Initialize(flagLogLevel); err != nil {
		panic(err)
	}

	ok := cfg.readStartParams()
	if !ok {
		flag.PrintDefaults()
		os.Exit(0)
	}
	collectMetrics()
}
