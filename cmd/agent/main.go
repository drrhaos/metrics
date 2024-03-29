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
const gaugesTotalMem = "TotalMemory"
const gaugesFreeMem = "FreeMemory"
const gaugesCPUutil = "CPUutilization1"

const urlUpdateMetricsJSONConst = "http://%s/updates/"

var sleepStep = map[uint]int64{0: 1, 1: 3, 2: 5}

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
